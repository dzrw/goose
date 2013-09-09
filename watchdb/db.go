package watchdb

import (
	"fmt"
)

type watchdb struct {
	index   map[string]*IntSet
	m       map[int]*Watch
	last_id int
}

func (db *watchdb) Size() int {
	return len(db.m)
}

// Clears the database.
func (db *watchdb) Clear() {
	db.m = make(map[int]*Watch)
	db.index = make(map[string]*IntSet)
	return
}

func (db *watchdb) Contains(id int) (ok bool) {
	if id >= 1 {
		_, ok = db.m[id]
		return
	}

	return false
}

// Removes a watch from the database.  Returns true if the watch was removed;
// otherwise, false if the watch was not found.
func (db *watchdb) Remove(id int) (ok bool) {
	if id >= 1 {
		if w, ok := db.m[id]; ok {
			delete(db.m, w.id)

			key := matchKey(w.Path, w.Method)
			db.index[key].Remove(w.id)
			return true
		}
	}

	return false
}

// Adds a watch to the database, returning an id and true if it was added;
// otherwise, false if an existing watch was overwritten.
func (db *watchdb) Add(w *Watch) (id int, ok bool) {
	key := matchKey(w.Path, w.Method)
	if _, ok := db.index[key]; !ok {
		db.index[key] = NewIntSet()
	}

	id, found := db.scan(w)
	if !found {
		db.last_id += 1
		id = db.last_id
	}

	// autogenerate the tag unless the tag was specified
	if w.Tag == "" {
		w.Tag = fmt.Sprintf("goose:events:%d", id)
	}

	w.id = id
	db.m[id] = w
	ok = db.index[key].Add(id)

	return
}

// Searches the database for any watches matching the specified request.
func (db *watchdb) Match(expr *MatchExpr) (m MatchData, ok bool) {
	key := matchKey(expr.Path, expr.Method)
	if index, ok := db.index[key]; ok {
		for id, _ := range index.set {
			m = NewMatchData(db.m[id])
			return m, true // whatever, loop once. todo. fix. etc.
		}

		// Index contained an empty set.
		return EmptyMatchData, false
	} else {
		// Not found in the index.
		return EmptyMatchData, false
	}
}

func matchKey(path, method string) string {
	return fmt.Sprintf("\"%s\"|\"%s\"", path, method)
}

// Scans the index for the specified watch.
func (db *watchdb) scan(w *Watch) (id int, ok bool) {
	key := matchKey(w.Path, w.Method)
	if idx, ok := db.index[key]; ok {
		for id, _ := range idx.set {
			wi := db.m[id]
			if wi.Path == w.Path && wi.Method == w.Method {
				return id, true
			}
		}
	}

	return
}
