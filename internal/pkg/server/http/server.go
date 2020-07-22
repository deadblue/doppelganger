package http

import "github.com/deadblue/doppelganger/internal/pkg/engine"

type Server struct {
	e *engine.Engine
}

func New(eng *engine.Engine) *Server {
	return &Server{
		e: eng,
	}
}
