package web

import (
	"fmt"
	"github.com/politician/goose/watchdb"
	"github.com/politician/goose/worker"
	"log"
	"net/http"
)

const (
	CRAZY_ERROR   = 509
	CRAZIER_ERROR = 510
)

// Intercepts incoming HTTP requests, and echoes prearranged
// responses.
type mallory struct {
	DataSourceName string
	m              worker.Matcher
}

func NewMatchHandler(name string, m worker.Matcher) *mallory {
	return &mallory{name, m}
}

// Implements http.Handler
func (d *mallory) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// The standard library does not support keep-alive timeouts,
	// so we'll just have to close the connection after every request
	// (or we could hope that clients remember to close their
	// connections (but if they don't, then we can't shutdown)).
	w.Header().Set("Connection", "close")

	m, err := d.lookup(req)
	if err != nil {
		http.Error(w, err.Error(), CRAZY_ERROR)
		return
	}

	d.render(w, m)
}

func (d *mallory) lookup(req *http.Request) (m watchdb.MatchData, err error) {
	expr := watchdb.NewMatchExpr(req.URL.Path, req.Method)
	m, err = d.m.Match(expr)

	str := ""
	if !m.IsMatch() {
		str = "not "
	}

	log.Printf("*** %s %s (%smatched)", req.Method, req.URL.Path, str)
	return
}

// Writes a response to the client on a match.
func (d *mallory) render(w http.ResponseWriter, m watchdb.MatchData) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}

		var str string

		err, ok := e.(error)
		if ok {
			str = err.Error()
		} else {
			str = fmt.Sprintf("%+v", e)
		}

		log.Printf("panic: (%d) %s", CRAZIER_ERROR, str)
		http.Error(w, str, CRAZIER_ERROR)
	}()

	w.Header().Set("X-Service-For", d.DataSourceName)
	w.Header().Set("X-Match", fmt.Sprintf("%t", m.IsMatch()))
	w.Header().Set("X-Match-Watch-ID", fmt.Sprintf("%d", m.Id()))

	if !m.IsMatch() {
		http.Error(w, "This request did not match any watches.", CRAZY_ERROR)
		return
	}

	m.Echo().ServeHTTP(w, nil)
}
