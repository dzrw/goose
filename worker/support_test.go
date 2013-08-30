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
	return watchdb.NewWatch("GET", "/foo", "redis:watch:1234", new(fakeResponse))
}

type fakeEventProvider struct{}

func (*fakeEventProvider) Dial() error {
	return nil
}

func (*fakeEventProvider) Close() {}

func assertSizeEquals(t *testing.T, mgr Manager, expected int) {
	status, err := mgr.Status()
	if err != nil {
		t.Error(err)
	}

	if actual := status.Size(); actual != expected {
		t.Errorf("expected: %d, got: %d", expected, actual)
	}
}
