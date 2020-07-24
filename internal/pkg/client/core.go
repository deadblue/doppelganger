package client

import (
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"github.com/deadblue/doppelganger/internal/pkg/server/raw"
	"github.com/deadblue/doppelganger/internal/pkg/suid"
	"log"
)

func (c *Client) call(method string, params []byte) (err error) {
	req := &pb.Request{
		Id:      suid.Next().String(),
		Method:  method,
		Request: params,
	}
	log.Printf("Request ID: %s, method: %s", req.Id, req.Method)
	err = raw.WriteMessage(c.conn, req)
	return
}

func (c *Client) Close() error {
	return c.conn.Close()
}
