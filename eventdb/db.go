package eventdb

import (
	"github.com/garyburd/redigo/redis"
	_ "log"
)

type eventdb struct {
	net  string
	addr string
	db   int
	conn redis.Conn
}

func (db *eventdb) Dial() (err error) {
	// Connect to redis.
	// dialer := NewRedisDialer(db.net, db.addr, 0)
	// conn, err := dialer.Dial()
	// if err != nil {
	// 	log.Fatalln("could not establish connection to redis: ", err)
	// }

	// db.conn = conn
	return
}

func (db *eventdb) Close() {
	// db.conn.Close()
}
