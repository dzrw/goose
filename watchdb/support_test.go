package watchdb

import (
	"net/http"
)

type fakeResponse struct{}

func (*fakeResponse) Status() int {
	return 200
}

func (*fakeResponse) Headers() map[string]string {
	return make(map[string]string)
}

func (*fakeResponse) Body() []byte {
	return make([]byte, 16)
}

func (*fakeResponse) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	panic("not implemented")
}

func fakeWatch() *Watch {
	return NewWatch("GET", "/foo", "redis:watch:1234", new(fakeResponse))
}
