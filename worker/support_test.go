package worker

import (
	"github.com/politician/goose/watchdb"
	"net/http"
	"testing"
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

func fakeWatch() *watchdb.Watch {
	return watchdb.NewWatch("dummy-service", "GET", "/foo", "redis:watch:1234", new(fakeResponse))
}

func assertSizeEquals(t *testing.T, ws WatchService, expected int) {
	status, err := ws.Status()
	if err != nil {
		t.Error(err)
	}

	if actual := status.Size(); actual != expected {
		t.Errorf("expected: %d, got: %d", expected, actual)
	}
}
