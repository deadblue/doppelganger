package task

import (
	"bytes"
	"context"
	"github.com/deadblue/gostream/quietly"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpTask struct {
	baseTask
	url     string
	method  string
	headers map[string]string
	body    []byte
}

func (t *HttpTask) Run(ctx context.Context) (err error) {
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
	if err != nil {
		return
	}
	defer quietly.Close(resp.Body)
	// Read response and send to callback
	if t.cb != nil {
		result, _ := ioutil.ReadAll(resp.Body)
		go t.cb.Send(result)
	} else {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
	}
	return nil
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
