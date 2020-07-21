package engine

import (
	"container/list"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type _QueueExecutor struct {
	// Wait group for waiting all tasks done
	wg *sync.WaitGroup
	// Closed flag
	closed int32

	// Blocking task channel between producer and consumer
	tasks chan Task
	// Infinite task queue
	queue *list.List
	// Task queue size
	qsize int32
	// New task event
	nt chan struct{}

	// Completed task count
	complete int
	// Total time for task
	duration time.Duration
}

// Submit add a task to executor's task queue.
func (e *_QueueExecutor) Submit(task Task) error {
	if e.closed != 0 {
		return errExecutorClosed
	}

	e.queue.PushBack(task)
	newSize := atomic.AddInt32(&e.qsize, 1)
	if newSize == 1 {
		e.nt <- struct{}{}
	}
	return nil
}

/*
Shutdown notifies executor to close.

After that, caller can not submit new task to this executor, but all running and
waiting tasks will be run.
*/
func (e *_QueueExecutor) Shutdown() {
	// Close only once
	if atomic.SwapInt32(&e.closed, 1) == 0 {
		// Notify producer to close.
		close(e.nt)
	}
}

// start goroutines for producer and consumers.
func (e *_QueueExecutor) start(n int) *_QueueExecutor {
	// Start consumers
	for i := 0; i < n; i++ {
		go e.consumer()
	}
	// Start producer
	go e.producer()
	return e
}

// producer carries task from queue to channel, and wait for new task.
func (e *_QueueExecutor) producer() {
	for alive := true; alive; {
		for node := e.queue.Front(); node != nil; node = e.queue.Front() {
			value := e.queue.Remove(node)
			atomic.AddInt32(&e.qsize, -1)
			node = nil // dereference
			if task, ok := value.(Task); ok {
				e.tasks <- task
			} else {
				// What the hell?
			}
		}
		// Waiting for new task or shutdown
		select {
		case _, ok := <-e.nt:
			//
			if !ok {
				alive = false
			}
		}
	}
	// No more task, close channel.
	close(e.tasks)
}

// consumer consumes task one by one.
func (e *_QueueExecutor) consumer() {
	e.wg.Add(1)
	defer e.wg.Done()

	for task := range e.tasks {
		// TODO: Retry when necessary.
		err := task.Run(context.Background())
		if err != nil {
			log.Printf("Run task error: %s", err)
		}
	}
}

// newQE creates an queue executor.
// It will fail if wait group is nil or core size is invalid.
func newQE(wg *sync.WaitGroup, opts *QueueOpts) *_QueueExecutor {
	if wg == nil || opts == nil || opts.CoreSize <= 0 {
		return nil
	}
	return (&_QueueExecutor{
		wg:     wg,
		closed: 0,
		tasks:  make(chan Task, opts.CoreSize),
		queue:  list.New(),
		qsize:  0,
		nt:     make(chan struct{}, 1),
	}).start(opts.CoreSize)
}
