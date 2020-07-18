package engine

import "sync"

type Engine struct {
	// wg for waiting all task workers done
	wg     *sync.WaitGroup
	queues map[string]*Queue
}

func (e *Engine) Submit(queue string, task TaskSpec) {
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
