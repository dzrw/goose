package eventdb

type EventProvider interface {
	Dial() error
	Close()
}

func NopEventProvider() EventProvider {
	return &nopdb{}
}

func RedisEventProvider(net, addr string, db int) EventProvider {
	return &redisdb{net, addr, db, nil}
}
