package eventdb

import (
	"net/http"
)

type EventProvider interface {
	Dial() error
	Close()
	Submit(tag string, req *http.Request)
}

func NopEventProvider() EventProvider {
	return &nopdb{}
}

func RedisEventProvider(net, addr string, db int) EventProvider {
	return &redisdb{net, addr, db, nil}
}
