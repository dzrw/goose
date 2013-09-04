package eventdb

import (
	"log"
	"net/http"
)

type nopdb struct{}

func (*nopdb) Dial() error {
	return nil
}

func (*nopdb) Close() {
	return
}

func (*nopdb) Trace(tag string, req *http.Request) {
	log.Printf("*** %s %s (matched)", req.Method, req.URL.Path)
	return
}

func (*nopdb) TraceUnexpected(req *http.Request) {
	log.Printf("*** %s %s (not matched)", req.Method, req.URL.Path)
	return
}
