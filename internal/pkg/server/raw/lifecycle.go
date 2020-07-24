package raw

import (
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"github.com/deadblue/doppelganger/internal/pkg/protocol"
	"github.com/deadblue/gostream/quietly"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"sync/atomic"
)

func (s *Server) startup() {
	log.Printf("Raw server listening at: %s", s.l.Addr())
	for alive := true; alive; {
		conn, err := s.l.Accept()
		if err != nil {
			if atomic.LoadInt32(&s.closing) == 1 {
				alive = false
			} else {
				log.Printf("Unexpected accept error: %s", err)
			}
		} else {
			go s.serve(conn)
		}
	}
}

func (s *Server) serve(conn net.Conn) {
	s.wg.Add(1)
	defer func() {
		quietly.Close(conn)
		s.wg.Done()
	}()
	// Read message one by one
	for atomic.LoadInt32(&s.closing) != 1 {
		req := &pb.Request{}
		if err := protocol.ReadMessage(conn, req); err != nil {
			if err != io.EOF {
				log.Printf("Read incoming message error: %s", err)
			}
			break
		}
		// Handle request
		resp := pb.Response{
			Id: req.Id, Error: 0,
		}
		if result, err := s.handle(req); err != nil {
			resp.Error = 1
			resp.Message = err.Error()
		} else {
			if result != nil {
				resp.Result, _ = proto.Marshal(result)
			}
		}
		if err := protocol.WriteMessage(conn, &resp); err != nil {
			log.Printf("Write outgoing message error: %s", err)
			break
		}
	}
}

func (s *Server) shutdown() {
	// Close listener
	quietly.Close(s.l)
	// Wait for all connections exit
	s.wg.Wait()
	// Close channel
	close(s.doneCh)
}
