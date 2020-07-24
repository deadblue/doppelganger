package raw

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

/*
For raw socket server, the messages between client and server are transfered
as stream, so we need define a procotol to carry the message by which the client
and server can explicitly pick messages from the stream.

The protocol is simple, each protocol package has following parts:

  * size: 2 bytes, big-endian unsigned int16 value.
  * body: up to 65535 bytes.

*/

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
