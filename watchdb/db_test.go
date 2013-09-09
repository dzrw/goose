package watchdb

import (
	"fmt"
	"testing"
)

func exampleWatch() *Watch {
	return NewWatch("dummy-service",
		"GET", "/foo", "redis:watch:1234",
		new(fakeResponse))
}

func checkAdd(t *testing.T, wp WatchProvider, w *Watch) (id int, ok bool) {
	id, ok = wp.Add(w)
	if !ok {
		t.Error("Could not add the watch.")
		return
	}

	if !wp.Contains(id) {
		t.Error("The database does not contain the watch ID.")
		return
	}

	return
}

func TestAdd(t *testing.T) {
	checkAdd(t, New(), exampleWatch())
}

func TestAddWithoutTag(t *testing.T) {
	wp := New()

	w := exampleWatch()
	w.Tag = ""

	id, ok := checkAdd(t, wp, w)
	if !ok {
		return
	}

	expected := fmt.Sprintf("goose:events:%d", id)
	if w.Tag != expected {
		t.Errorf("expected tag=%s, got tag=%s", expected, w.Tag)
		return
	}
}

func TestMatch(t *testing.T) {
	wp := New()
	w := exampleWatch()

	id, ok := checkAdd(t, wp, w)
	if !ok {
		return
	}

	mx := NewMatchExpr(w.DataSourceName, w.Path, w.Method)
	m, ok := wp.Match(mx)

	switch {
	case !ok:
		t.Error("Expected a match, but none were found")
		return

	case m.Id() != id:
		t.Errorf("expected id=%s, got id=%s", id, m.Id())
		return
	}
}

func TestMatchWhenDataSourceNameIsDifferent(t *testing.T) {
	wp := New()
	w := exampleWatch()

	_, ok := checkAdd(t, wp, w)
	if !ok {
		return
	}

	mx := NewMatchExpr("wrong", w.Path, w.Method)
	if _, ok := wp.Match(mx); ok {
		t.Errorf("expected not to match on DSN (wdsn=%s, mdsn=%s)", w.DataSourceName, mx.DataSourceName)
		return
	}
}

func TestRemoveWhenNotExists(t *testing.T) {
	wp := New()
	if ok := wp.Remove(42); ok {
		t.Errorf("expected %t, got %t", false, ok)
		return
	}
}

func TestRemoveWhenExists(t *testing.T) {
	wp := New()
	w := exampleWatch()

	id, ok := checkAdd(t, wp, w)
	if !ok {
		return
	}

	if ok := wp.Remove(id); !ok {
		t.Errorf("expected %t, got %t", true, ok)
		return
	}

	if wp.Contains(id) {
		t.Errorf("The db still contains a watch with id=%d", id)
		return
	}
}

func TestMatchAfterRemove(t *testing.T) {
	wp := New()
	w := exampleWatch()

	id, _ := wp.Add(w)
	wp.Remove(id)

	mx := NewMatchExpr(w.DataSourceName, w.Path, w.Method)
	_, ok := wp.Match(mx)
	if ok {
		t.Errorf("expected %t, got %t", false, ok)
		return
	}
}
