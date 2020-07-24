/*
package protocol provides function for read/write messages between
raw socket client and server.

The protocol is simple, each protocol package has following parts:

  * size: 2 bytes, big-endian unsigned int16 value.
  * body: up to 65535 bytes.

*/
package protocol
