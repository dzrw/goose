package web

import (
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/politician/goose/watchdb"
	"github.com/politician/goose/worker"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMatchHandler(t *testing.T) {
	const (
		LISTEN_ADDR = ":9002"
		WATCHED_URL = "http://127.0.0.1:9002/foo"
		WATCH_BODY  = "hello, world"
	)

	mgr, err := worker.Start(nil)
	if err != nil {
		t.Error(err)
		return
	}

	defer mgr.Stop()

	w := fakeWatch("redis:watch:1234", "GET", "/foo", WATCH_BODY)
	_, err = mgr.Add(w)
	if err != nil {
		t.Error(err)
		return
	}

	m, _ := mgr.Match(watchdb.NewMatchExpr("/foo", "GET"))
	if !m.IsMatch() {
		t.Errorf("watchdb doesn't recognize the watch")
		return
	}

	mux := NewMatchHandler("test", mgr)

	s, err := StartHttpServer(LISTEN_ADDR, mux)
	if err != nil {
		t.Error(err)
		return
	}

	defer s.Stop()

	resp, err := http.Get(WATCHED_URL)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected: %d, got: %d", 200, resp.StatusCode)
		return
	}

	bufr := bufio.NewReader(resp.Body)
	yy, err := ioutil.ReadAll(bufr)
	resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}

	str := string(yy)
	if str != WATCH_BODY {
		t.Errorf("unexpected: %s", str)
		t.Errorf("expected: \"%s\", got: \"%s\"", WATCH_BODY, str)

		// Fprintln adds a '\n', oops.
		spew.Dump(yy)
		spew.Dump(str)
		return
	}
}

type fakeResponse struct {
	body string
}

func (*fakeResponse) Status() int {
	return 200
}

func (*fakeResponse) Headers() map[string]string {
	return make(map[string]string)
}

func (f *fakeResponse) Body() []byte {
	return []byte(f.body)
}

func (f *fakeResponse) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, f.body)
}

func fakeWatch(tag, method, path, body string) *watchdb.Watch {
	echo := &fakeResponse{body: body}
	return watchdb.NewWatch(path, method, tag, echo)
}
