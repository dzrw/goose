package eventdb

type EventProvider interface {
	Dial() error
	Close()
}

func Dial(net, addr string) (db EventProvider, err error) {
	db = &eventdb{net: net, addr: addr}
	err = db.Dial()
	return db, err
}
