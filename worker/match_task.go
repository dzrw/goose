package worker

import (
	"github.com/politician/goose/watchdb"
)

type matchTask struct {
	expr *watchdb.MatchExpr
	ch   chan watchdb.MatchData
}

func (t *matchTask) await() (v watchdb.MatchData, err error) {
	v, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

func (t *matchTask) resolve(v watchdb.MatchData) {
	t.ch <- v
}
