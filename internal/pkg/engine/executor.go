package engine

import (
	"container/list"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type _Job struct {
	t  Task
	cb Callback
}

type _Executor struct {
	// Wait group for waiting all tasks done
	wg *sync.WaitGroup
	// Closed flag
	cf int32

	// Job queue
	jq *list.List
	// Job queue size
	js int32
	// Blocking task channel between producer and consumer
	jch chan _Job
	// Event bus for new job and shutdown
	eb chan struct{}

	// Completed task count
	complete int
	// Total time for task
	duration time.Duration
}

// Submit add a task to executor's task queue.
func (e *_Executor) Submit(task Task, cb Callback) error {
	if e.cf != 0 {
		return errExecutorClosed
	}
	e.jq.PushBack(_Job{
		t:  task,
		cb: cb,
	})
	newSize := atomic.AddInt32(&e.js, 1)
	if newSize == 1 {
		e.eb <- struct{}{}
	}
	return nil
}

/*
Shutdown notifies executor to close.

After that, caller can not submit new task to this executor, but all running and
waiting tasks will be run.
*/
func (e *_Executor) Shutdown() {
	// Close only once
	if atomic.SwapInt32(&e.cf, 1) == 0 {
		// Notify producer to close.
		close(e.eb)
	}
}

// start goroutines for producer and consumers.
func (e *_Executor) start(n int) *_Executor {
	// Start consumers
	for i := 0; i < n; i++ {
		go e.consumer()
	}
	// Start producer
	go e.producer()
	return e
}

// producer carries task from queue to channel, and wait for new task.
func (e *_Executor) producer() {
	for alive := true; alive; {
		for node := e.jq.Front(); node != nil; node = e.jq.Front() {
			value := e.jq.Remove(node)
			atomic.AddInt32(&e.js, -1)
			node = nil // dereference
			if job, ok := value.(_Job); ok {
				e.jch <- job
			} else {
				// What the hell?
			}
		}
		// Waiting for new task or shutdown
		select {
		case _, ok := <-e.eb:
			if !ok {
				// executor has been shutdown
				alive = false
			}
		}
	}
	// No more task, close channel.
	close(e.jch)
}

// consumer consumes task one by one.
func (e *_Executor) consumer() {
	e.wg.Add(1)
	defer e.wg.Done()

	for job := range e.jch {
		// TODO: Retry and send result to callback.
		err := job.t.Run(context.Background())
		if err != nil {
			log.Printf("Run task error: %s", err)
		} else {
			e.complete += 1
		}
	}
}

// newExecutor creates an queue executor.
// It will fail if wait group is nil or core size is invalid.
func newExecutor(wg *sync.WaitGroup, opts *QueueOpts) *_Executor {
	if wg == nil || opts == nil || opts.CoreSize <= 0 {
		return nil
	}
	return (&_Executor{
		wg: wg,
		cf: 0,

		jq:  list.New(),
		js:  0,
		eb:  make(chan struct{}, 1),
		jch: make(chan _Job, opts.CoreSize),
	}).start(opts.CoreSize)
}
