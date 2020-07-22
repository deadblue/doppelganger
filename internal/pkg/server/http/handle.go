package http

import (
	"net/http"
)

const (
	headerContentType = "Content-Type"

	ctJson     = "application/json"
	ctProtobuf = "application/protobuf"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: handle http request.
}
