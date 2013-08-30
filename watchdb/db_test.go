package watchdb

import (
	"testing"
)

func TestAdd(t *testing.T) {
	db := New()
	w := NewWatch("GET", "/foo", "redis:watch:1234", new(fakeResponse))

	id, ok := db.Add(w)
	if !ok {
		t.Error("Things are not ok. OK? Leave me alone. /wrists")
	}

	if !db.Contains(id) {
		t.Error("The database does not contain the watch ID.")
	}
}

func TestMatch(t *testing.T) {
	db := New()
	w := NewWatch("GET", "/foo", "redis:watch:1234", new(fakeResponse))

	id, _ := db.Add(w)
	m, ok := db.Match(NewMatchExpr("GET", "/foo"))

	switch {
	case !ok:
		t.Error("Expected a match, but none were found")

	case m.Id() != id:
		t.Errorf("expected id=%s, got id=%s", id, m.Id())
	}
}

func TestRemoveWhenNotExists(t *testing.T) {
	db := New()
	if ok := db.Remove(42); ok {
		t.Errorf("expected %t, got %t", false, ok)
	}
}

func TestRemoveWhenExists(t *testing.T) {
	db := New()
	w := NewWatch("GET", "/foo", "redis:watch:1234", new(fakeResponse))

	id, ok := db.Add(w)
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	if ok := db.Remove(id); !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	if db.Contains(id) {
		t.Errorf("The db still contains a watch with id=%d", id)
	}
}

func TestMatchAfterRemove(t *testing.T) {
	db := New()
	w := NewWatch("GET", "/foo", "redis:watch:1234", new(fakeResponse))

	id, _ := db.Add(w)
	db.Remove(id)

	if _, ok := db.Match(NewMatchExpr(w.Path, w.Method)); ok {
		t.Errorf("expected %t, got %t", false, ok)
	}
}
