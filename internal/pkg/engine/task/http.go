package task

import (
	"bytes"
	"context"
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"github.com/deadblue/gostream/quietly"
	"io"
	"net/http"
)

type HttpTask struct {
	url     string
	method  string
	headers map[string]string
	body    []byte
}

func (t *HttpTask) Run(ctx context.Context, cb engine.Callback) (err error) {
	// Make request
	if t.method == "" {
		if t.body == nil {
			t.method = http.MethodGet
		} else {
			t.method = http.MethodPost
		}
	}
	body := io.Reader(nil)
	if t.body != nil {
		body = bytes.NewReader(t.body)
	}
	req, err := http.NewRequestWithContext(ctx, t.method, t.url, body)
	if err != nil {
		return
	}
	// Set headers
	req.Header.Set("User-Agent", "Doppelganger/1.0")
	if t.headers != nil {
		for name, value := range t.headers {
			req.Header.Set(name, value)
		}
	}
	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer quietly.Close(resp.Body)
		if cb != nil {
			err = cb.Send(resp.Body)
		}
	}
	return
}

func (t *HttpTask) Header(name, value string) *HttpTask {
	t.headers[name] = value
	return t
}

func HttpGet(url string) *HttpTask {
	return Http(url, http.MethodGet, nil)
}

func HttpPost(url string, data []byte) *HttpTask {
	return Http(url, http.MethodPost, data)
}

func Http(url, method string, data []byte) *HttpTask {
	task := &HttpTask{
		url:     url,
		method:  method,
		headers: map[string]string{},
	}
	if data != nil {
		task.body = make([]byte, len(data))
		copy(task.body, data)
	}
	return task
}
