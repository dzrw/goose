package worker

import (
	"errors"
	// "github.com/davecgh/go-spew/spew"
	"github.com/politician/goose/watchdb"
	"log"
)

// WatchKeeper serializes access to an underlying watch
// database and Redis endpoint.  If contention around this
// point becomes a problem, we have several options including
// using a buffered (instead of unbufferred) channel, or
// building a pool of these objects.  Or do both.  It's not
// going to be a big problem.
type worker struct {
	db        watchdb.WatchProvider
	work      chan Task
	sigs      SignalChannel
	verbosity int
}

var ErrWatchIdOutOfRange = errors.New("argument must be greater than 0")

// Spawns a goroutine to manage concurrent access to the watch database
// and redis backend.
func (self *worker) Start() {
	go func() { self.forever() }()
	return
}

// Convenience method to instruct the backend to gracefully shutdown.
func (self *worker) Stop() int {
	if self.verbosity > 0 {
		log.Println("goose/worker: watch service stopping")
	}
	exitCode := self.sigs.Quit()
	return exitCode
}

// Convenience method to issue a probe (query) to the backend.
func (self *worker) Match(expr *watchdb.MatchExpr) (m watchdb.MatchData, err error) {
	tr, err := self.schedule(NewMatchTask(expr))
	if err != nil {
		return
	}

	m, err = UnpackMatchTaskResponse(tr)
	return
}

// Convenience method to clear the watch database.
func (self *worker) Clear() (err error) {
	_, err = self.remove(FLUSH_DB)
	return
}

// Convenience method to remove a watch from the watch database.
func (self *worker) Remove(id int) (removed bool, err error) {
	if id <= 0 {
		return false, ErrWatchIdOutOfRange
	}

	removed, err = self.remove(id)
	return
}

// Schedules a remove watch task.
func (self *worker) remove(id int) (removed bool, err error) {
	tr, err := self.schedule(NewRemoveTask(id))
	if err != nil {
		return
	}

	removed, err = UnpackRemoveTaskResponse(tr)
	return
}

// Convenience method to add a watch to the watch database.
func (self *worker) Add(w *watchdb.Watch) (id int, tag string, err error) {
	tr, err := self.schedule(NewCreateTask(w))
	if err != nil {
		return
	}

	id, tag, err = UnpackCreateTaskResponse(tr)
	return
}

// Returns information about the state of the worker.
func (self *worker) Status() (status Status, err error) {
	tr, err := self.schedule(NewStatusTask())
	if err != nil {
		return
	}

	status, err = UnpackStatusTaskResponse(tr)
	return
}

func (self *worker) UnsafeGetWatchProvider() watchdb.WatchProvider {
	return self.db
}

// Runs tasks in the background until signalled.
func (self *worker) forever() {
	for {
		select {
		// Interrupts arrive on the sigs channel.
		case s := <-self.sigs:
			if s == nil {
				if self.verbosity > 0 {
					log.Println("goose/worker: supervisor channel closed")
				}
				return
			}

			switch s.Signal() {
			case SIGQUIT:
				if self.verbosity > 0 {
					log.Println("goose/worker: watch service stopped")
				}
				s.Resolve(1)
				return
			}

		// Normal work arrives on the work channel.
		case t := <-self.work:
			err := t.Do(self.db)
			if err != nil {
				log.Panicf("%T task failed: %+v", t, err)
			}
		}
	}
}

// Schedules a task onto the background goroutine.
func (self *worker) schedule(t Task) (tr TaskResponse, err error) {
	self.work <- t
	tr, err = t.Await()
	return
}
