package raw

import (
	"encoding/binary"
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"github.com/deadblue/gostream/quietly"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"sync/atomic"
	"time"
)

func (s *Server) startup() {
	log.Printf("Raw server listening at: %s", s.l.Addr())
	for {
		conn, err := s.l.Accept()
		if err != nil {
			log.Printf("Accept error: %s", err)
			if atomic.LoadInt32(&s.closed) == 1 {
				break
			}
		} else {
			go s.serve(&clientConn{
				Conn:   conn,
				active: false,
			})
		}
	}
}

func (s *Server) serve(conn *clientConn) {
	s.conns[conn] = struct{}{}
	defer func() {
		delete(s.conns, conn)
		quietly.Close(conn)
	}()
	for {
		req := &pb.Request{}
		if err := readMessage(conn, req); err == nil {
			conn.active = true
			s.handle(req)
			conn.active = false
		} else {
			log.Printf("Read incoming message error: %s", err)
			break
		}
	}
}

func readMessage(r io.Reader, req *pb.Request) (err error) {
	// Read message size
	buf := make([]byte, 2)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	size := binary.BigEndian.Uint16(buf)
	log.Printf("Message size: %d", size)
	// Read raw message
	buf = make([]byte, size)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	return proto.Unmarshal(buf, req)
}

func (s *Server) shutdown() {
	// Close listener, DO NOT accept new connection.
	quietly.Close(s.l)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for ok := false; !ok; {
		for conn := range s.conns {
			if !conn.active {
				quietly.Close(conn)
			}
		}
		ok = len(s.conns) == 0
		select {
		case <-ticker.C:
		}
	}

	// close channel
	close(s.doneCh)
}
