package task

import "github.com/deadblue/doppelganger/internal/pkg/engine"

type CallbackTask interface {
	engine.Task
	Callback(cb engine.Callback)
}

type baseTask struct {
	cb engine.Callback
}

func (t *baseTask) Callback(cb engine.Callback) {
	t.cb = cb
}
