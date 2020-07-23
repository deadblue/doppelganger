package http

import (
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"net"
	"net/http"
)

type Server struct {
	// Core engine
	e *engine.Engine
	// Http server
	hs *http.Server

	// Close flag
	closed int32
	// Error channel
	errCh chan error
}

func (s *Server) Shutdown() {
	go s.shutdown()
}

func (s *Server) Error() <-chan error {
	return s.errCh
}

func New(e *engine.Engine, l net.Listener) *Server {
	s := &Server{
		e:      e,
		hs:     &http.Server{},
		closed: 0,
		errCh:  make(chan error),
	}
	s.hs.Handler = s
	go s.startup(l)
	return s
}
