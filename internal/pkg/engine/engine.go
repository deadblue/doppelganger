package engine

import (
	"sync"
)

type Engine struct {
	// wg for waiting all task workers done
	wg *sync.WaitGroup
	// Task queue
	queues map[string]*queue
}

func (e *Engine) Submit(queue string, task Task) {
	if q, ok := e.queues[queue]; ok {
		q.Submit(task)
	} else {

	}
}

func (e *Engine) Wait() {
	e.wg.Wait()
}

func New() *Engine {
	return &Engine{
		wg: &sync.WaitGroup{},
	}
}
