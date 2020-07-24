package json

import (
	"encoding/json"
)

type CommandCallback struct {
	// Command name
	Name string `json:"name"`
	// Command args
	Args []string `json:"args"`
}

type HttpCallback struct {
	// Config URL
	Url string `json:"url"`
	// Request headers
	Headers map[string]string `json:"headers"`
}

type Callback struct {
	// Callback type: command/http
	Type CallbackType `json:"type"`
	// Callback target.
	Config json.RawMessage `json:"config"`
}

type CommandTask struct {
	CommandCallback
	// Input data
	Input string `json:"input"`
}

type HttpTask struct {
	HttpCallback
	// Request method
	Method string `json:"method"`
	// Request body
	Body string `json:"body"`
}

type TaskParams struct {
	// Target queue name
	Queue string `json:"queue"`
	// Task type: command/http
	Type TaskType `json:"type"`
	// Task config
	Config json.RawMessage `json:"config"`
	// Maximum retry times
	RetryTimes int `json:"retry_times"`
	// Callback config.
	Callback *Callback `json:"callback"`
}
