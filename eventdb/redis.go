package eventdb

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisDialer struct {
	Network string
	Address string
	Db      int
}

func NewRedisDialer(network, addr string, db int) *RedisDialer {
	if network == "" {
		network = "tcp"
	}

	return &RedisDialer{network, addr, db}
}

func (d *RedisDialer) Dial() (conn redis.Conn, err error) {
	conn, err = redis.Dial(d.Network, d.Address)
	if err != nil {
		return
	}

	log.Printf("connected to redis on %s:%s\n", d.Network, d.Address)

	// TODO Select the db.

	return
}
