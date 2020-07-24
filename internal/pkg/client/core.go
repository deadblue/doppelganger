package client

import (
	"errors"
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"github.com/deadblue/doppelganger/internal/pkg/protocol"
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
	if err = protocol.WriteMessage(c.conn, req); err != nil {
		return
	}
	resp := &pb.Response{}
	if err = protocol.ReadMessage(c.conn, resp); err == nil {
		if resp.Error != 0 {
			err = errors.New(resp.Message)
		} else {
			log.Printf("Result: %s", resp.Result)
		}
	}
	return
}
