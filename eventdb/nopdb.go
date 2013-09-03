package eventdb

import (
	"net/http"
)

type nopdb struct{}

func (*nopdb) Dial() error {
	return nil
}

func (*nopdb) Close() {
	return
}

func (*nopdb) Submit(tag string, req *http.Request) {
	return
}
