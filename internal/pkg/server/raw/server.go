package raw

import (
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"net"
	"sync/atomic"
)

type clientConn struct {
	net.Conn
	active bool
}

type Server struct {
	// core engine
	e *engine.Engine
	// network listener
	l net.Listener
	// connections holder
	conns map[*clientConn]struct{}
	// closed flag
	closed int32
	// shutdown channel
	doneCh chan struct{}
}

func (s *Server) Shutdown() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		go s.shutdown()
	}
}

func (s *Server) Done() <-chan struct{} {
	return s.doneCh
}

func New(e *engine.Engine, l net.Listener) *Server {
	s := &Server{
		e:      e,
		l:      l,
		conns:  make(map[*clientConn]struct{}),
		doneCh: make(chan struct{}),
	}
	go s.startup()
	return s
}
