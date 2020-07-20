package callback

import (
	"bytes"
	"github.com/deadblue/gostream/quietly"
	"log"
	"net/http"
)

type HttpCallback struct {
	url     string
	headers map[string]string
}

func (c *HttpCallback) Send(result []byte) {
	// Make request
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(result))
	if err != nil {
		log.Printf("Create callback request failed: %s", err)
		return
	}
	if c.headers != nil {
		for name, value := range c.headers {
			req.Header.Set(name, value)
		}
	}
	// Send request
	if resp, err := http.DefaultClient.Do(req); err != nil {
		log.Printf("Send callback request error: %s", err)
	} else {
		defer quietly.Close(resp.Body)
		log.Printf("Callback request status code: %d", resp.StatusCode)
	}
}

func (c *HttpCallback) Header(name, value string) *HttpCallback {
	c.headers[name] = value
	return c
}

// Http creates a HTTP request based callback.
func Http(url string) *HttpCallback {
	return &HttpCallback{
		url:     url,
		headers: map[string]string{},
	}
}
