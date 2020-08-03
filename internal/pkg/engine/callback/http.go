package callback

import (
	"github.com/deadblue/gostream/quietly"
	"io"
	"log"
	"net/http"
)

type HttpCallback struct {
	// Callback URL
	url string
	// Callback HTTP headers
	hdrs map[string]string
}

func (hc *HttpCallback) Send(r io.Reader) (err error) {
	// Make request
	req, err := http.NewRequest(http.MethodPost, hc.url, r)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Doppelganger/1.0")
	if hc.hdrs != nil {
		for name, value := range hc.hdrs {
			req.Header.Set(name, value)
		}
	}
	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("HTTP callback error: %s", err)
	} else {
		defer quietly.Close(resp.Body)
		log.Printf("HTTP callback status code: %d", resp.StatusCode)
	}
	return
}

func (hc *HttpCallback) Header(name, value string) *HttpCallback {
	hc.hdrs[name] = value
	return hc
}

func Http(url string) *HttpCallback {
	return &HttpCallback{
		url:  url,
		hdrs: make(map[string]string),
	}
}
