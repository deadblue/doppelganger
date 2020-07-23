package http

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync/atomic"
)

func (s *Server) startup(l net.Listener) {
	log.Printf("HTTP server listened at: %s", l.Addr())
	if err := s.hs.Serve(l); err != http.ErrServerClosed {
		s.close(err)
	}
}

func (s *Server) shutdown() {
	if atomic.LoadInt32(&s.closed) == 0 {
		err := s.hs.Shutdown(context.Background())
		s.close(err)
	}
}

func (s *Server) close(err error) {
	// Close only once
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		if err != nil {
			s.errCh <- err
		}
		close(s.errCh)
	}
}
