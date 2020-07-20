package engine

import (
	"context"
	"log"
	"sync"
	"time"
)

type queue struct {
	wg *sync.WaitGroup

	tasks chan Task

	// Completed task count
	complete int
	// Total time for task
	duration time.Duration
}

func (q *queue) Submit(task Task) {
	q.tasks <- task
}

func (q *queue) Shutdown() {
	close(q.tasks)
}

func (q *queue) StartWorkers(n int) {
	for i := 0; i < n; i++ {
		go q.worker()
	}
}

func (q *queue) worker() {
	q.wg.Add(1)
	defer q.wg.Done()

	for task := range q.tasks {
		err := task.Run(context.Background())
		if err != nil {
			log.Printf("Run task error: %s", err)
		}
	}
}

func newQueue(wg *sync.WaitGroup, size int) *queue {
	if size <= 0 {
		return nil
	}

	q := &queue{
		wg:    wg,
		tasks: make(chan Task, 1000),
	}
	q.StartWorkers(size)
	return q
}
