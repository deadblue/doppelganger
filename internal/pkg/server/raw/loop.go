package raw

import (
	"encoding/binary"
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"github.com/deadblue/doppelganger/internal/pkg/suid"
	"github.com/deadblue/gostream/quietly"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
)

func (s *Server) accept() {
	for alive := true; alive; {
		conn, err := s.l.Accept()
		if err != nil {
			alive = false
		} else {
			go s.serve(suid.Next(), conn)
		}
	}
}

func (s *Server) serve(id string, conn net.Conn) {
	log.Printf("New connection: #%s", id)
	s.cm[id] = conn
	defer func() {
		log.Printf("Closing connection: %s", id)
		quietly.Close(conn)
		delete(s.cm, id)
	}()

	for {
		// Read message size
		buf := make([]byte, 2)
		if _, err := io.ReadFull(conn, buf); err != nil {
			if err != io.EOF {
				log.Printf("Unexpected error when read size: %s", err)
			}
			break
		}
		size := binary.BigEndian.Uint16(buf)
		// Read raw message
		buf = make([]byte, size)
		if _, err := io.ReadFull(conn, buf); err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				log.Printf("Unexpected error when read message: %s", err)
			}
			break
		}
		// Parse incoming message
		req := &pb.Request{}
		if err := proto.Unmarshal(buf, req); err == nil {
			s.handle(req)
		} else {
			log.Printf("Parse incoming message error: %s", err)
			break
		}
	}
}

func (s *Server) handle(req *pb.Request) {
	log.Printf("Handle request: %s", req.Method)
	// TODO: dispatch the request.
}
