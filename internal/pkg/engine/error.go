package engine

import "errors"

var (
	// When caller submit task with a non-exists queue.
	errUnknownQueue = errors.New("unknown queue")
	// When caller submit task to a closed queue executor.
	errExecutorClosed = errors.New("task exector has been closed")
)
