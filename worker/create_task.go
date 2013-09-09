package worker

import (
	"github.com/politician/goose/watchdb"
)

type createTask struct {
	w  *watchdb.Watch
	ch chan *createTaskResponse
}

type createTaskResponse struct {
	id  int
	tag string
}

func NewCreateTask(w *watchdb.Watch) *createTask {
	return &createTask{w, make(chan *createTaskResponse)}
}

func (t *createTask) Do(repo watchdb.WatchProvider) (err error) {
	id, ok := repo.Add(t.w)
	if !ok {
		err = ErrCannotAddWatch
		return
	}

	t.ch <- &createTaskResponse{id, t.w.Tag}
	return
}

func (t *createTask) Await() (tr TaskResponse, err error) {
	tr, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

// Unpacks a CreateTask's TaskResponse into an actual response.
func UnpackCreateTaskResponse(tr TaskResponse) (id int, tag string, err error) {
	res, ok := tr.(*createTaskResponse)
	if !ok {
		err = ErrWrongType
		return
	}

	id = res.id
	tag = res.tag
	return
}
