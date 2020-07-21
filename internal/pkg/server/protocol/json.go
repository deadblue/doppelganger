package protocol

import (
	"encoding/json"
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"github.com/deadblue/doppelganger/internal/pkg/engine/callback"
	"github.com/deadblue/doppelganger/internal/pkg/engine/task"
)

type CallbackCommandJson struct {
	// Command name
	Name string `json:"name"`
	// Command args
	Args []string `json:"args"`
}

type TaskCommandJson struct {
	CallbackCommandJson
	// Input data
	Input []byte `json:"input"`
}

type CallbackHttpJson struct {
	// Config URL
	Url string `json:"url"`
	// Request headers
	Headers map[string]string `json:"headers"`
}

type TaskHttpJson struct {
	CallbackHttpJson
	// Request method
	Method string `json:"method"`
	// Request body
	Body []byte `json:"body"`
	// Maximum retry times
	RetryTimes int `json:"retry_times"`
}

type CallbackJson struct {
	// Callback type: command/http
	Type CallbackType `json:"type"`
	// Callback target.
	Config json.RawMessage `json:"config"`
}

type TaskJson struct {
	// Task type: command/http
	Type TaskType `json:"type"`
	// Task config
	Config json.RawMessage `json:"config"`
	// Maximum retry times
	RetryTimes int `json:"retry_times"`
	// Callback config.
	Callback *CallbackJson `json:"callback"`
}

func (j *TaskJson) Parse() (t engine.Task, err error) {
	// Parse task
	ct := task.CallbackTask(nil)
	if j.Type.IsCommand() {
		config := &TaskCommandJson{}
		if err = json.Unmarshal(j.Config, config); err != nil {
			return
		}
		cmdTask := task.Command(config.Name, config.Args)
		if config.Input != nil {
			cmdTask.Input(config.Input)
		}
		ct = cmdTask
	} else if j.Type.IsHttp() {
		config := &TaskHttpJson{}
		if err = json.Unmarshal(j.Config, config); err != nil {
			return
		}
		httpTask := task.Http(config.Url, config.Method, config.Body)
		if config.Headers != nil {
			for name, value := range config.Headers {
				httpTask.Header(name, value)
			}
		}
		ct = httpTask
	} else {
		err = errUnknownTaskType
		return
	}
	// Parse task callback
	if cb := j.Callback; cb != nil {
		if cb.Type.IsCommand() {
			config := &CallbackCommandJson{}
			if err = json.Unmarshal(cb.Config, config); err == nil {
				ct.Callback(callback.Command(config.Name, config.Args))
			}
		} else if cb.Type.IsHttp() {
			config := &CallbackHttpJson{}
			if err = json.Unmarshal(cb.Config, config); err == nil {
				hc := callback.Http(config.Url)
				if config.Headers != nil {
					for name, value := range config.Headers {
						hc.Header(name, value)
					}
				}
				ct.Callback(hc)
			}
		} else {
			err = errUnknownCallbackType
		}
		if err != nil {
			t = nil
			return
		}
	}
	t = ct
	return
}
