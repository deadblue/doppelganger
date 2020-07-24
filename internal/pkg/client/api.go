package client

func (c *Client) Greeting(name string) (err error) {
	return c.call("greet", []byte(name))
}
