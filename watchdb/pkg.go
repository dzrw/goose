package watchdb

import (
	"net/http"
)

type WatchProvider interface {
	Size() (n int)
	Clear()
	Remove(id int) (ok bool)
	Contains(id int) (ok bool)
	Add(w *Watch) (id int, ok bool)
	Match(expr *MatchExpr) (m MatchData, ok bool)
}

// Initializes a new database.
func New() (db WatchProvider) {
	db = &watchdb{}
	db.Clear()
	return
}

type Echo interface {
	http.Handler

	Status() int
	Headers() map[string]string
	Body() []byte
}

type MatchData interface {
	Id() int
	Tag() string
	Echo() Echo
	IsMatch() bool
}
