package engine

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

type SimpleTask struct {
	index int
}

func (t *SimpleTask) Run(_ context.Context) error {
	log.Printf("Hello, world => %d", t.index)
	time.Sleep(1 * time.Second)
	return nil
}

func TestExecutor(t *testing.T) {
	wg := &sync.WaitGroup{}
	e := newQE(wg, 3)
	for i := 1; i <= 50; i++ {
		log.Printf("Submit task => %d", i)
		_ = e.Submit(&SimpleTask{index: i})
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("All task submit, shutting down ...")
	e.Shutdown()

	log.Println("Waiting for task done!")
	wg.Wait()
	log.Println("All task done!")
}
