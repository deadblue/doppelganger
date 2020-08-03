package raw

import (
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Server struct {
	// core engine
	e *engine.Engine
	// network listener
	l net.Listener

	// closing flag
	cf int32
	// wait group for all connections closed
	wg sync.WaitGroup
	// channel to notify server has completely shutdown.
	doneCh chan struct{}
}

func (s *Server) Shutdown() {
	if atomic.CompareAndSwapInt32(&s.cf, 0, 1) {
		go s.shutdown()
	} else {
		log.Println("Server is being / has been shutdown!")
	}
}

func (s *Server) Done() <-chan struct{} {
	return s.doneCh
}

func New(e *engine.Engine, l net.Listener) *Server {
	s := &Server{
		// task engine
		e: e,
		// network listener
		l: l,
		// closing flag
		cf: 0,
		// wait for all active connections exit.
		wg: sync.WaitGroup{},
		// notification channel
		doneCh: make(chan struct{}),
	}
	go s.startup()
	return s
}
