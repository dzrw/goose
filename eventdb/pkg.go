package eventdb

type EventProvider interface {
	Dial() error
	Close()

	Trace(m *Message) error
}

func NopEventProvider() EventProvider {
	return &nopdb{}
}

func RedisEventProvider(net, addr string, db int, verbosity int) EventProvider {
	return &redisdb{net, addr, db, nil, verbosity}
}
