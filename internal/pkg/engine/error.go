package engine

import "errors"

var (
	errQueueNotExist = errors.New("queue not exist")

	errQueueExist = errors.New("queue already exist")

	errExecutorClosed = errors.New("task exector has been closed")
)
