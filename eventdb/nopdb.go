package eventdb

import (
	"log"
)

type nopdb struct{}

func (*nopdb) Dial() error {
	return nil
}

func (*nopdb) Close() {
	return
}

func (*nopdb) Trace(m *Message) error {
	log.Print(m.String())
	return nil
}
