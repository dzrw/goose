package worker

import (
	"github.com/politician/goose/watchdb"
)

type createTask struct {
	w  *watchdb.Watch
	ch chan int
}

func (t *createTask) await() (v int, err error) {
	v, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

func (t *createTask) resolve(v int) {
	t.ch <- v
}
