package protocol

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

func ReadMessage(r io.Reader, msg proto.Message) (err error) {
	// Read message size
	buf := make([]byte, 2)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	size := binary.BigEndian.Uint16(buf)
	// Read message data
	buf = make([]byte, size)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	return proto.Unmarshal(buf, msg)
}
