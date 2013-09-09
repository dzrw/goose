package worker

import (
	"github.com/politician/goose/watchdb"
)

type WatchService interface {
	Stop() int
	Start()
	Status() (status Status, err error)

	Add(w *watchdb.Watch) (id int, tag string, err error)
	Remove(id int) (removed bool, err error)
	Clear() (err error)

	Match(expr *watchdb.MatchExpr) (m watchdb.MatchData, err error)
}

func Start() (ws WatchService, err error) {
	ws = &worker{
		db:   watchdb.New(),
		work: make(chan Task),
		sigs: make(SignalChannel),
	}

	ws.Start()
	return
}
