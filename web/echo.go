package web

import (
	"github.com/politician/goose/watchdb"
	"net/http"
)

func NewEcho(status int, headers map[string]string, body []byte) watchdb.Echo {
	return &echo{
		status:  status,
		headers: headers,
		body:    body,
	}
}

var NotFound = NewEcho(http.StatusNotFound, make(map[string]string), make([]byte, 0))

type echo struct {
	status  int
	headers map[string]string
	body    []byte
}

func (res *echo) Status() int {
	return res.status
}

func (res *echo) Headers() map[string]string {
	return res.headers
}

func (res *echo) Body() []byte {
	return res.body
}

func (res *echo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for k, v := range res.Headers() {
		w.Header().Set(k, v)
	}

	w.WriteHeader(res.Status())

	_, err := w.Write(res.Body())
	if err != nil {
		panic(err)
	}
}
