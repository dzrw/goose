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

	d.m.ServeHTTP(w, req)
}
