package client

import (
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"github.com/deadblue/doppelganger/internal/pkg/server/raw"
	"github.com/deadblue/doppelganger/internal/pkg/suid"
	"google.golang.org/protobuf/proto"
	"log"
)

func (c *Client) call(method string, params proto.Message) (err error) {
	req := &pb.Request{
		Id:     suid.Next().String(),
		Method: method,
	}
	if params != nil {
		if data, err := proto.Marshal(params); err == nil {
			req.Params = data
		} else {
			return err
		}
	}
	log.Printf("Request ID: %s, method: %s", req.Id, req.Method)
	err = raw.WriteMessage(c.conn, req)
	return
}
