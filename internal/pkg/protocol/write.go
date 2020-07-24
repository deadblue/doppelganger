package protocol

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

func WriteMessage(w io.Writer, msg proto.Message) (err error) {
	// Marshal message
	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	// Write message size
	sizeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(sizeBuf, uint16(len(buf)))
	if _, err = w.Write(sizeBuf); err != nil {
		return
	}
	// Write message data
	_, err = w.Write(buf)
	return
}
