package core

import (
	"sync"
	"time"
)

type Queue struct {
	wg *sync.WaitGroup

	tasks chan interface{}

	complete int
	time     time.Duration
}

func (q *Queue) Submit(task interface{}) {
	q.tasks <- task
}

func (q *Queue) Shutdown() {
	close(q.tasks)
}
