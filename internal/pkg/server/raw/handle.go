package raw

import (
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"log"
)

func (s *Server) handle(req *pb.Request) {
	// TODO: handle the request
	log.Printf("Handle request: %s", req.Method)
}
