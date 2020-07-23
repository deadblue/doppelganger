package http

import (
	"log"
	"net/http"
)

const (
	headerContentType = "Content-Type"

	ctJson     = "application/json"
	ctProtobuf = "application/protobuf"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: handle http request.
	log.Printf("Request uri => %s", r.URL.RequestURI())
	w.Header().Set("X-Powered-By", "Doppelganger")
	w.WriteHeader(http.StatusOK)
}
