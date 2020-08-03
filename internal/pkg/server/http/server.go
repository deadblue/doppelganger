package http

import (
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"net"
	"net/http"
)

type Server struct {
	// http server
	hs *http.Server
	// close flag
	cf int32
	// error channel
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
		hs: &http.Server{
			Handler: &handler{e: e},
		},
		cf:    0,
		errCh: make(chan error),
	}
	go s.startup(l)
	return s
}
