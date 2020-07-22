package engine

import (
	"log"
	"sync"
)

type Engine struct {
	// WaitGroup for waiting all tasks done, it will be shared
	// across all queue executors.
	wg *sync.WaitGroup

	// Mutex for adding queue.
	qm *sync.Mutex
	// Task queue executors.
	qes map[string]*_QueueExecutor

	// Close event channel
	done chan struct{}
}

func (e *Engine) startUp() {
	e.wg.Add(1)
	e.wg.Wait()
	e.done <- struct{}{}
	close(e.done)
}

func (e *Engine) Queue(name string, opts *QueueOpts) {
	e.qm.Lock()
	defer e.qm.Unlock()

	if _, ok := e.qes[name]; !ok {
		e.qes[name] = newQE(e.wg, opts)
	}
}

func (e *Engine) Submit(name string, task Task) (err error) {
	if qe, ok := e.qes[name]; ok {
		err = qe.Submit(task)
	} else {
		err = errUnknownQueue
	}
	return
}

func (e *Engine) Shutdown() {
	for n, qe := range e.qes {
		log.Printf("Shutting down queue [%s]", n)
		qe.Shutdown()
	}
	e.Done()
}

func (e *Engine) Done() <-chan struct{} {
	return e.done
}

func New() (engine *Engine) {
	checkUser()
	engine = &Engine{
		wg:   &sync.WaitGroup{},
		qes:  make(map[string]*_QueueExecutor),
		done: make(chan struct{}),
	}
	go engine.startUp()
	return
}
