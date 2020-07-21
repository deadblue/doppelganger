package protocol

import "errors"

var (
	//
	errUnknownTaskType = errors.New("unknown task type")
	//
	errUnknownCallbackType = errors.New("unknown callback type")
)
