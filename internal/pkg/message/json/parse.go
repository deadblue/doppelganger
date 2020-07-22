package json

import (
	"encoding/json"
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"github.com/deadblue/doppelganger/internal/pkg/engine/callback"
	"github.com/deadblue/doppelganger/internal/pkg/engine/task"
)

func (r *TaskRequest) Parse() (t engine.Task, err error) {
	// Parse task
	ct := task.CallbackTask(nil)
	if r.Type.IsCommand() {
		config := &CommandTask{}
		if err = json.Unmarshal(r.Config, config); err != nil {
			return
		}
		cmdTask := task.Command(config.Name, config.Args)
		if config.Input != "" {
			cmdTask.Input([]byte(config.Input))
		}
		ct = cmdTask
	} else if r.Type.IsHttp() {
		config := &HttpTask{}
		if err = json.Unmarshal(r.Config, config); err != nil {
			return
		}
		httpTask := task.Http(config.Url, config.Method, []byte(config.Body))
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
	if cb := r.Callback; cb != nil {
		if cb.Type.IsCommand() {
			config := &CommandCallback{}
			if err = json.Unmarshal(cb.Config, config); err == nil {
				ct.Callback(callback.Command(config.Name, config.Args))
			}
		} else if cb.Type.IsHttp() {
			config := &HttpCallback{}
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
