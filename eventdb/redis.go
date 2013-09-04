package eventdb

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

var ErrCouldNotSelectDb = errors.New("could not select db")

type redisdb struct {
	net  string
	addr string
	db   int
	conn redis.Conn
}

func (db *redisdb) Dial() (err error) {
	conn, err := redis.Dial(db.net, db.addr)
	if err != nil {
		return
	}

	db.conn = conn

	_, err = redis.String(db.conn.Do("SELECT", db.db))
	if err != nil {
		db.conn.Close()
		return
	}

	log.Printf("connected to redis on %s:%s\n", db.net, db.addr)
	return
}

func (db *redisdb) Close() {
	db.conn.Close()
}

func (db *redisdb) Trace(tag string, req *http.Request) {
	log.Printf("*** %s %s (matched)", req.Method, req.URL.Path)
	return
}

func (db *redisdb) TraceUnexpected(req *http.Request) {
	log.Printf("*** %s %s (not matched)", req.Method, req.URL.Path)
	return
}
