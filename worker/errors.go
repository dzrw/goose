package worker

import (
	"errors"
)

var ErrWrongType = errors.New("wrong type")
var ErrCannotRemoveWatchZero = errors.New("cannot remove watch 0")
var ErrChannelClosed = errors.New("the channel was closed unexpectedly")
var ErrCannotAddWatch = errors.New("cannot add the watch")
