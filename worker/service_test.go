package worker

import (
	"github.com/politician/goose/watchdb"
	_ "log"
	"testing"
)

func TestStartAndStopServiceTwice(t *testing.T) {
	mgr := Start(&fakeEventProvider{})

	if exitCode := mgr.Stop(); exitCode != 1 {
		t.Errorf("expected: %d, got: %d", 1, exitCode)
	}

	mgr.Start()
	mgr.Stop()
}

func TestStatusWithService(t *testing.T) {
	mgr := Start(&fakeEventProvider{})
	defer mgr.Stop()

	assertSizeEquals(t, mgr, 0)
}

func TestAddWithService(t *testing.T) {
	mgr := Start(&fakeEventProvider{})
	defer mgr.Stop()

	w := fakeWatch()

	_, err := mgr.Add(w)
	if err != nil {
		t.Error(err)
	}

	assertSizeEquals(t, mgr, 1)
}

func TestClearWithService(t *testing.T) {
	mgr := Start(&fakeEventProvider{})
	defer mgr.Stop()

	w := fakeWatch()

	_, err := mgr.Add(w)
	if err != nil {
		t.Error(err)
	}

	assertSizeEquals(t, mgr, 1)

	if err := mgr.Clear(); err != nil {
		t.Error(err)
	}

	assertSizeEquals(t, mgr, 0)
}

func TestMatchWithService(t *testing.T) {
	mgr := Start(&fakeEventProvider{})
	defer mgr.Stop()

	w := fakeWatch()
	mgr.Add(w)

	expr := watchdb.NewMatchExpr(w.Path, w.Method)
	m, err := mgr.Match(expr)
	if err != nil {
		t.Error(err)
	}

	if m.Tag() != w.Tag {
		t.Error("expected: %s, got: %s", w.Tag, m.Tag())
	}
}
