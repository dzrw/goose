package worker

import (
	"errors"
	// "github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/politician/goose/eventdb"
	"github.com/politician/goose/watchdb"
	"log"
	"net/http"
)

// WatchKeeper serializes access to an underlying watch
// database and Redis endpoint.  If contention around this
// point becomes a problem, we have several options including
// using a buffered (instead of unbufferred) channel, or
// building a pool of these objects.  Or do both.  It's not
// going to be a big problem.
type worker struct {
	db       watchdb.WatchProvider
	provider eventdb.EventProvider
	work     chan T
	sigs     SignalChannel
}

type T interface{}

var ErrWatchIdOutOfRange = errors.New("argument must be greater than 0")

// Spawns a goroutine to manage concurrent access to the watch database
// and redis backend.
func (wk *worker) Start() {
	go func() { wk.forever() }()
	return
}

// Convenience method to instruct the backend to gracefully shutdown.
func (wk *worker) Stop() int {
	log.Println("worker: stopping")
	exitCode := wk.sigs.Quit()
	wk.provider.Close()
	return exitCode
}

// Convenience method to issue a probe (query) to the backend.
func (wk *worker) Match(expr *watchdb.MatchExpr) (m watchdb.MatchData, err error) {
	t := matchTask{expr, make(chan watchdb.MatchData)}
	wk.enqueue(t)
	m, err = t.await()
	return
}

// Convenience method to clear the watch database.
func (wk *worker) Clear() (err error) {
	t := removeTask{FLUSH_DB, make(chan bool)}
	wk.enqueue(t)
	_, err = t.await()
	return
}

// Convenience method to remove a watch from the watch database.
func (wk *worker) Remove(id int) (removed bool, err error) {
	if id <= 0 {
		return false, ErrWatchIdOutOfRange
	}

	t := removeTask{id, make(chan bool)}
	wk.enqueue(t)
	removed, err = t.await()
	return
}

// Convenience method to add a watch to the watch database.
func (wk *worker) Add(w *watchdb.Watch) (id int, err error) {
	t := createTask{w, make(chan int)}
	wk.enqueue(t)
	id, err = t.await()
	return
}

// Returns information about the state of the worker.
func (wrk *worker) Status() (status Status, err error) {
	t := statusTask{make(chan Status)}
	wrk.enqueue(t)
	status, err = t.await()
	return
}

func (wrk *worker) UnsafeGetWatchProvider() watchdb.WatchProvider {
	return wrk.db
}

func (wrk *worker) forever() {
	for {
		select {
		// Interrupts arrive on the sigs channel.
		case s := <-wrk.sigs:
			if s == nil {
				log.Println("worker: supervisor channel closed")
				return
			}

			switch s.Signal() {
			case SIGQUIT:
				log.Println("worker: stopped")
				s.Resolve(1)
				return
			}

		// Normal work arrives on the work channel.
		case t := <-wrk.work:
			switch u := t.(type) {
			case matchTask:
				wrk.match(u)
			case createTask:
				wrk.add(u)
			case removeTask:
				wrk.remove(u)
			case statusTask:
				wrk.status(u)
			default:
				log.Fatalf("unexpected task: %T %+v", t, t)
			}
		}
	}
}

func (wk *worker) enqueue(t T) {
	wk.work <- t
}

func (wrk *worker) match(t matchTask) {
	m, _ := wrk.db.Match(t.expr)
	t.resolve(m)
}

func (wrk *worker) add(t createTask) {
	id, _ := wrk.db.Add(t.w)

	t.resolve(id)
}

func (wrk *worker) remove(t removeTask) {
	if t.id == FLUSH_DB {
		wrk.db.Clear()
		t.resolve(true)
	} else if ok := wrk.db.Remove(t.id); ok {
		t.resolve(true)
	} else {
		t.resolve(false)
	}
}

func (wrk *worker) status(t statusTask) {
	size := wrk.db.Size()
	t.resolve(size)
}

// ----------------------------------------------
// Watch Resolver
// ----------------------------------------------

func (wrk *worker) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	const (
		GOOSE_ERROR      = 510
		GOOSE_CONF_ERROR = 509
	)

	defer func() {
		e := recover()
		if e == nil {
			return
		}

		var str string

		switch u := e.(type) {
		case error:
			str = u.Error()
		default:
			str = fmt.Sprintf("%+v", e)
		}

		log.Printf("panic: (%d) %s", GOOSE_ERROR, str)
		http.Error(w, str, GOOSE_ERROR)
	}()

	// Build a match expression from the request
	expr := watchdb.NewMatchExpr(req.URL.Path, req.Method)

	// Lookup the associated watch, if any
	m, err := wrk.Match(expr)
	if err != nil {
		http.Error(w, err.Error(), GOOSE_ERROR)
		return
	}

	// TODO: m.MatchType == watchdb.PASSTHRU

	// Log misses
	if !m.IsMatch() {
		wrk.provider.TraceUnexpected(req)
		http.Error(w, "This request did not match any watches.", GOOSE_CONF_ERROR)
		return
	}

	// Track hits, and serve echoes.
	wrk.provider.Trace(m.Tag(), req)
	m.Echo().ServeHTTP(w, req)
}
