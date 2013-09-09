package eventdb

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"time"
)

var ErrCouldNotSelectDb = errors.New("could not select db")

const (
	GOOSE_REDIS_UFO_LIST   = "goose:unexpected"
	GOOSE_REDIS_REQ_PREFIX = "goose:requests:"
)

type redisdb struct {
	net     string
	addr    string
	db      int
	pool    *redis.Pool
	verbose int
}

func (db *redisdb) Dial() (err error) {
	db.pool = makePool(db.net, db.addr, db.db)

	err = db.ping()
	if err != nil {
		return
	}

	if db.verbose > 0 {
		log.Printf("goose/eventdb: connected to redis on %s:%s\n", db.net, db.addr)
	}
	return
}

func (db *redisdb) ping() error {
	c := db.pool.Get()
	defer c.Close()
	_, err := c.Do("PING")
	return err
}

func (db *redisdb) Close() {
	db.pool.Close()

	if db.verbose > 0 {
		log.Println("goose/eventdb: redis connection pool released")
	}
}

func (db *redisdb) Trace(m *Message) (err error) {
	log.Print(m.String())

	// Ensure that we've got a tag.
	tag := m.Tag
	if !m.Tagged() {
		tag = GOOSE_REDIS_UFO_LIST
	}

	// Convert the request body to a base64-encoded string.
	input, err := ioutil.ReadAll(m.Request.Body)
	if err != nil {
		return
	}

	body, err := EncodeBuffer(input)
	if err != nil {
		return
	}

	// Merge message key/value pairs.
	args := make(map[string]string)
	copymap(args, m.MoreInfo)
	copymap(args, map[string]string{
		"tag":        tag,
		"timestamp":  string(m.Timestamp),
		"remoteAddr": m.Request.RemoteAddr,
		"method":     m.Request.Method,
		"url":        m.Request.URL.String(),
		"body":       body,
	})

	// Deliver the message to Redis.
	err = deliver(db.pool, tag, args)
	return
}

func deliver(pool *redis.Pool, key string, data map[string]string) (err error) {
	hkey, err := RandomKey(GOOSE_REDIS_REQ_PREFIX, 16)
	if err != nil {
		return
	}

	conn := pool.Get()
	defer conn.Close()

	// Record the request/match.
	for field, val := range data {
		if err = conn.Send("HSET", hkey, field, val); err != nil {
			return
		}
	}

	// Notify any processes blocking on the associated list.
	if err = conn.Send("LPUSH", key, hkey); err != nil {
		return
	}

	// Flush the pipeline.
	if err = conn.Flush(); err != nil {
		return
	}

	// Read all of the replies, but just drop them on the ground.
	for i := 0; i < len(data)+1; i += 1 {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}

	return
}

func makePool(net, addr string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (c redis.Conn, err error) {
			c, err = redis.Dial(net, addr)
			if err != nil {
				return
			}

			_, err = c.Do("SELECT", db)
			if err != nil {
				c.Close()
				return
			}

			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func copymap(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}
