package eventdb

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
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

	err := db.storeRequest(tag, req)
	if err != nil {
		panic(err)
	}
}

func (db *redisdb) TraceUnexpected(req *http.Request) {
	log.Printf("*** %s %s (not matched)", req.Method, req.URL.Path)

	err := db.storeRequest("goose:unexpected", req)
	if err != nil {
		panic(err)
	}
}

func (db *redisdb) storeRequest(list string, req *http.Request) (err error) {
	key, err := RandomKey("goose:requests:", 16)
	if err != nil {
		return
	}

	method := req.Method
	url := req.URL.String()
	body, err := RequestBodyToString(req)
	if err != nil {
		return
	}

	_, err = db.conn.Do("HMSET", key, "method", method, "url", url, "body", body)
	if err != nil {
		return
	}

	_, err = db.conn.Do("LPUSH", list, key)
	if err != nil {
		return
	}

	return
}

func RequestBodyToString(req *http.Request) (encodedBody string, err error) {
	input, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	str, err := EncodeBuffer(input)
	if err != nil {
		return
	}

	return str, err
}
