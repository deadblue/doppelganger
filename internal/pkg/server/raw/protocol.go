package raw

/*
For raw socket server, the messages between client and server are transfered
as stream, so we need define a procotol to carry the message by which the client
and server can explicitly pick messages from the stream.

The protocol is simple, each protocol package has following parts:

  * size: 2 bytes, big-endian unsigned int16 value.
  * body: up to 65535 bytes.

*/
