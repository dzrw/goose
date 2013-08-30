package worker

const (
	FLUSH_DB = -1
)

type removeTask struct {
	id int
	ch chan bool
}

func (t *removeTask) await() (v bool, err error) {
	v, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

func (t *removeTask) resolve(removed bool) (err error) {
	t.ch <- removed
	return
}
