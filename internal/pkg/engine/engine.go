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
	qm sync.Mutex
	// Job executors.
	es map[string]*_Executor
	// Close event channel
	done chan struct{}
}

func (e *Engine) startUp() {
	e.wg.Add(1)
	e.wg.Wait()
	e.done <- struct{}{}
	close(e.done)
}

// QueueAdd creates a queue in engine with specified name and opts.
func (e *Engine) QueueAdd(name string, opts *QueueOpts) (err error) {
	e.qm.Lock()
	defer e.qm.Unlock()

	if _, ok := e.es[name]; ok {
		err = errQueueExist
	} else {
		e.es[name] = newExecutor(e.wg, opts)
	}
	return
}

func (e *Engine) QueueList() {
	// TODO
	return
}

// JobAdd adds a task to queue.
// When the task done, its result will be sent.
func (e *Engine) JobAdd(queue string, task Task, callback Callback) (err error) {
	if ex, ok := e.es[queue]; ok {
		err = ex.Submit(task, callback)
	} else {
		err = errQueueNotExist
	}
	return
}

// Shutdown notify all executors to shutdown.
func (e *Engine) Shutdown() {
	for n, qe := range e.es {
		log.Printf("Shutting down queue [%s]", n)
		qe.Shutdown()
	}
	e.wg.Done()
}

func (e *Engine) Done() <-chan struct{} {
	return e.done
}

func New() (engine *Engine) {
	engine = &Engine{
		wg:   &sync.WaitGroup{},
		qm:   sync.Mutex{},
		es:   make(map[string]*_Executor),
		done: make(chan struct{}),
	}
	go engine.startUp()
	return
}
