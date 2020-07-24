package client

import "net"

// Client for raw socket server.
type Client struct {
	// connection to the server
	conn net.Conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func New(network, address string) (c *Client, err error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return
	}
	c = &Client{conn: conn}
	return
}
