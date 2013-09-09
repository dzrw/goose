package web

import (
	"encoding/json"
	_ "github.com/davecgh/go-spew/spew"
	"github.com/politician/goose/watchdb"
	"github.com/politician/goose/worker"
	"net/http"
	"strings"
	"testing"
)

type UnsafeAccessor interface {
	UnsafeGetWatchProvider() watchdb.WatchProvider
}

func TestCreateWatch(t *testing.T) {
	ws, err := worker.Start()
	if err != nil {
		t.Error(err)
		return
	}

	defer ws.Stop()

	mux := NewWatchHandler(ws)

	s, err := StartHttpServer(ListenAddr, mux)
	if err != nil {
		t.Error(err)
		return
	}

	defer s.Stop()

	id, err := postWatch(t, "http://127.0.0.1:9001/watches")
	if err != nil {
		t.Error(err)
		return
	}

	const WATCH_ID = 1

	if id != WATCH_ID {
		t.Errorf("expected: %d, got: %d", WATCH_ID, id)
		return
	}

	db := ws.(UnsafeAccessor).UnsafeGetWatchProvider()
	if !db.Contains(WATCH_ID) {
		t.Errorf("expected the watchdb to contain a watch with id=%d", WATCH_ID)
		return
	}
}

func postWatch(t *testing.T, addr string) (id int, err error) {
	data := &JsonWatch{
		Tag:            "opaque-id-1234",
		DataSourceName: "twitter",
		MatchExpr: &JsonMatchExpr{
			Method: "GET",
			Path:   "/foo/bar",
		},
		Echo: &JsonEcho{
			Status:  200,
			Headers: make(map[string]string),
			Body:    "hello, world",
		},
	}

	buf, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
		return
	}

	postBody := strings.NewReader(string(buf))

	resp, err := http.Post(addr, "application/json", postBody)
	if err != nil {
		t.Error(err)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	var tmp struct {
		Id  int
		Tag string
	}
	err = decoder.Decode(&tmp)
	if err != nil {
		t.Logf("%+v", decoder)
		t.Error(err)
		return
	}

	resp.Body.Close()

	return tmp.Id, nil
}
