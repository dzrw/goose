package worker

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

func (t *statusTask) await() (v Status, err error) {
	v, ok := <-t.ch
	if !ok {
		err = ErrChannelClosed
	}
	return
}

func (t *statusTask) resolve(size int) {
	t.ch <- &status{size}
}
