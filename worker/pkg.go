package worker

import (
	"github.com/politician/goose/eventdb"
	"github.com/politician/goose/watchdb"
)

type Matcher interface {
	Match(expr *watchdb.MatchExpr) (m watchdb.MatchData, err error)
}

type WatchList interface {
	Add(w *watchdb.Watch) (id int, err error)
	Remove(id int) (removed bool, err error)
	Clear() (err error)
}

type Manager interface {
	Matcher
	WatchList

	Stop() int
	Start()
	Status() (status Status, err error)
}

func Start(ep eventdb.EventProvider) (wrk Manager, err error) {
	if ep == nil {
		ep = eventdb.NopEventProvider()
	}

	err = ep.Dial()
	if err != nil {
		return
	}

	wrk = &worker{
		db:       watchdb.New(),
		provider: ep,
		work:     make(chan T),
		sigs:     make(SignalChannel),
	}

	wrk.Start()
	return
}
