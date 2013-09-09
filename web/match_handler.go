package web

import (
	"fmt"
	"github.com/politician/goose/eventdb"
	"github.com/politician/goose/watchdb"
	"github.com/politician/goose/worker"
	"log"
	"net/http"
)

// Intercepts incoming HTTP requests, and echoes prearranged
// responses.
type matchHandler struct {
	DataSourceName string
	ws             worker.WatchService
	ep             eventdb.EventProvider
}

func NewMatchHandler(dsn string, ws worker.WatchService, ep eventdb.EventProvider) *matchHandler {
	return &matchHandler{dsn, ws, ep}
}

// An HTTP handler for matching incoming requests with existing watches.
func (self *matchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// The standard library does not support keep-alive timeouts,
	// so we'll just have to close the connection after every request
	// (or we could hope that clients remember to close their
	// connections (but if they don't, then we can't shutdown)).
	w.Header().Set("Connection", "close")

	defer func() {
		if e := recover(); e != nil {
			var str string

			switch u := e.(type) {
			case error:
				str = u.Error()
			default:
				str = fmt.Sprintf("%+v", e)
			}

			log.Printf("panic: (%d) %s", GOOSE_ERROR, str)
			http.Error(w, str, GOOSE_ERROR)
		}
	}()

	// Build a match expression from the request
	expr := watchdb.NewMatchExpr(self.DataSourceName, req.URL.Path, req.Method)

	// Lookup the associated watch, if any
	m, err := self.ws.Match(expr)
	if err != nil {
		http.Error(w, err.Error(), GOOSE_ERROR)
		return
	}

	// TODO: m.MatchType == watchdb.PASSTHRU

	tag := ""
	if m.IsMatch() {
		tag = m.Tag()
	}

	// Record the event.
	self.ep.Trace(eventdb.NewMessage(tag, req, map[string]string{
		"dataSourceName": self.DataSourceName,
	}))

	// Return an error code if we don't find any matching watches.
	if !m.IsMatch() {
		http.Error(w, "This request did not match any watches.", GOOSE_CONF_ERROR)
		return
	}

	// Otherwise, return an echo.
	m.Echo().ServeHTTP(w, req)
}
