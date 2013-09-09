package worker

import (
	"github.com/politician/goose/watchdb"
)

const (
	FLUSH_DB = -1
)

type removeTask struct {
	id int
	ch chan bool
}

func NewRemoveTask(id int) *removeTask {
	return &removeTask{id, make(chan bool)}
}

func (t *removeTask) Do(repo watchdb.WatchProvider) (err error) {
	var removed bool

	if t.id == FLUSH_DB {
		repo.Clear()
		removed = true
	} else if ok := repo.Remove(t.id); ok {
		removed = true
	}

	t.ch <- removed
	return
}

func (t *removeTask) Await() (tr TaskResponse, err error) {
	tr, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

// Unpacks a RemoveTask's TaskResponse into an actual response.
func UnpackRemoveTaskResponse(tr TaskResponse) (removed bool, err error) {
	removed, ok := tr.(bool)
	if !ok {
		err = ErrWrongType
		return
	}

	return
}
