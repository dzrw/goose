package worker

import (
	"github.com/politician/goose/watchdb"
)

type Task interface {
	Do(repo watchdb.WatchProvider) (err error)
	Await() (tr TaskResponse, err error)
}

type TaskResponse interface{}
