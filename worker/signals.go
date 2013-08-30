package worker

type sigval int

const (
	SIGQUIT sigval = iota
)

// goroutine signals
type signal struct {
	val  sigval
	done chan int
}

func (s *signal) Signal() sigval {
	return s.val
}

func (s *signal) Resolve(val int) {
	s.done <- val
}

// -------------------------------------

type SignalChannel chan *signal

func (ch *SignalChannel) Quit() int {
	done := make(chan int)
	sig := &signal{SIGQUIT, done}
	*ch <- sig
	ret := <-done
	return ret
}
