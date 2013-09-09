package worker

import (
	"github.com/politician/goose/watchdb"
)

type Status interface {
	Size() int
}

type status struct {
	size int
}

func (s *status) Size() int {
	return s.size
}

type statusTask struct {
	ch chan Status
}

func NewStatusTask() *statusTask {
	return &statusTask{make(chan Status)}
}

func (t *statusTask) Do(repo watchdb.WatchProvider) (err error) {
	size := repo.Size()
	t.ch <- &status{size}
	return
}

func (t *statusTask) Await() (tr TaskResponse, err error) {
	tr, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

// Unpacks a StatusTask's TaskResponse into an actual response.
func UnpackStatusTaskResponse(tr TaskResponse) (status Status, err error) {
	status, ok := tr.(Status)
	if !ok {
		err = ErrWrongType
		return
	}

	return
}
