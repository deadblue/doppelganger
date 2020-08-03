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
	duration int64
}

// Submit add a task into to executor's queue.
func (e *_Executor) Submit(task Task, callback Callback) error {
	if e.cf != 0 {
		return errExecutorClosed
	}
	e.jq.PushBack(_Job{
		t:  task,
		cb: callback,
	})
	newSize := atomic.AddInt32(&e.js, 1)
	if newSize == 1 {
		e.eb <- struct{}{}
	}
	return nil
}

/*
Shutdown notifies executor to close.

After that, caller can not submit new task to this executor, but all tasks
in queue will be done.
*/
func (e *_Executor) Shutdown() {
	// Close only once
	if atomic.SwapInt32(&e.cf, 1) == 0 {
		// Notify producer to close.
		close(e.eb)
	}
}

// start goroutines for producer and consumers.
func (e *_Executor) startup(n int) *_Executor {
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
	// No more task, close job channel.
	close(e.jch)
}

// consumer consumes task one by one.
func (e *_Executor) consumer() {
	e.wg.Add(1)
	defer e.wg.Done()

	for job := range e.jch {
		if err := e.doJob(&job); err != nil {
			log.Printf("Run job error: %s", err)
		}
	}
}

func (e *_Executor) doJob(job *_Job) (err error) {
	// TODO: how to retry?
	startTime := time.Now()
	err = job.t.Run(context.Background(), job.cb)
	runningTime := time.Now().Sub(startTime)
	if err != nil {
		log.Printf("Run task error: %s", err)
	} else {
		// Update statistics
		e.complete += 1
		atomic.AddInt64(&e.duration, int64(runningTime))
	}
	return
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
	}).startup(opts.CoreSize)
}
