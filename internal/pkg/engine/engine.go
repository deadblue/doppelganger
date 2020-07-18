package core

type Engine struct {
	queues map[string]*Queue
}

func (e *Engine) Submit(queue string, task TaskSpec) {
	if q, ok := e.queues[queue]; ok {
		q.Submit(task)
	} else {

	}
}
