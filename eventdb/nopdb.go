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

func (*nopdb) Trace(tag string, req *http.Request, res *http.Response) {
	return
}

func (*nopdb) TraceUnexpected(req *http.Request, res *http.Response) {
	return
}
