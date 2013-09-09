package worker

import (
	"github.com/politician/goose/watchdb"
	_ "log"
	"testing"
)

func TestStartAndStopServiceTwice(t *testing.T) {
	ws, err := Start()
	if err != nil {
		t.Error(err)
		return
	}

	if exitCode := ws.Stop(); exitCode != 1 {
		t.Errorf("expected: %d, got: %d", 1, exitCode)
		return
	}

	ws.Start()
	ws.Stop()
}

func TestStatusWithService(t *testing.T) {
	ws, err := Start()
	if err != nil {
		t.Error(err)
		return
	}

	defer ws.Stop()

	assertSizeEquals(t, ws, 0)
}

func TestAddWithService(t *testing.T) {
	ws, err := Start()
	if err != nil {
		t.Error(err)
		return
	}

	defer ws.Stop()

	w := fakeWatch()

	_, _, err = ws.Add(w)
	if err != nil {
		t.Error(err)
		return
	}

	assertSizeEquals(t, ws, 1)
}

func TestClearWithService(t *testing.T) {
	ws, err := Start()
	if err != nil {
		t.Error(err)
		return
	}

	defer ws.Stop()

	w := fakeWatch()

	_, _, err = ws.Add(w)
	if err != nil {
		t.Error(err)
		return
	}

	assertSizeEquals(t, ws, 1)

	if err := ws.Clear(); err != nil {
		t.Error(err)
		return
	}

	assertSizeEquals(t, ws, 0)
}

func TestMatchWithService(t *testing.T) {
	ws, err := Start()
	if err != nil {
		t.Error(err)
		return
	}

	defer ws.Stop()

	w := fakeWatch()
	ws.Add(w)

	expr := watchdb.NewMatchExpr("dummy-service", w.Path, w.Method)
	m, err := ws.Match(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if m.Tag() != w.Tag {
		t.Error("expected: %s, got: %s", w.Tag, m.Tag())
		return
	}
}
