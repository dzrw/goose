package worker

import (
	"github.com/politician/goose/watchdb"
)

type matchTask struct {
	expr *watchdb.MatchExpr
	ch   chan watchdb.MatchData
}

func NewMatchTask(expr *watchdb.MatchExpr) *matchTask {
	return &matchTask{expr, make(chan watchdb.MatchData)}
}

func (t *matchTask) Do(repo watchdb.WatchProvider) (err error) {
	m, _ := repo.Match(t.expr)
	t.ch <- m
	return
}

func (t *matchTask) Await() (tr TaskResponse, err error) {
	tr, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

// Unpacks a MatchTask's TaskResponse into an actual response.
func UnpackMatchTaskResponse(tr TaskResponse) (m watchdb.MatchData, err error) {
	m, ok := tr.(watchdb.MatchData)
	if !ok {
		err = ErrWrongType
		return
	}

	return
}
