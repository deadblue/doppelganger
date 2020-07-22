package raw

import (
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"github.com/deadblue/gostream/quietly"
	"log"
	"net"
)

type Server struct {
	// core engine
	e *engine.Engine
	// network listener
	l net.Listener
	// connect manager
	cm map[string]net.Conn
	// shutdown channel
	done chan struct{}
}

func (s *Server) Start() {
	log.Printf("Raw server listening at: %s", s.l.Addr())
	go s.accept()
}

func (s *Server) Shutdown() (err error) {
	// Close the listener
	err = s.l.Close()
	// Force close all connections
	for id, conn := range s.cm {
		log.Printf("Closing connection #%s", id)
		quietly.Close(conn)
	}
	close(s.done)
	return
}

func (s *Server) Done() <-chan struct{} {
	return s.done
}

func New(eng *engine.Engine) *Server {
	return &Server{
		e: eng,

		done: make(chan struct{}),
	}
}
