package raw

import (
	"errors"
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"github.com/deadblue/doppelganger/internal/pkg/engine/callback"
	"github.com/deadblue/doppelganger/internal/pkg/engine/task"
	"github.com/deadblue/doppelganger/internal/pkg/message"
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"google.golang.org/protobuf/proto"
	"log"
)

var (
	errUnknownMethod = errors.New("unknown method")

	errInvalidParams = errors.New("parameters not enough")
)

func (s *Server) handle(req *pb.Request) (result proto.Message, err error) {
	switch req.Method {
	case message.MethodTaskAdd:
		result, err = s.handleAddTask(req.Params)
	default:
		err = errUnknownMethod
	}
	return
}

func (s *Server) handleAddTask(paramsData []byte) (result proto.Message, err error) {
	params := &pb.AddTaskParams{}
	if err = proto.Unmarshal(paramsData, params); err != nil {
		return
	}
	if params.Queue == "" || params.Task == nil {
		err = errInvalidParams
		return
	}
	// Parse callback
	cb := engine.Callback(nil)
	if pcb := params.Callback; pcb != nil {
		switch pcb.Type {
		case pb.CallbackType_CB_COMMAND:
			config := pb.CommandCallback{}
			err = proto.Unmarshal(pcb.Config, &config)
			if err == nil {
				cb = callback.Command(config.Name, config.Args...)
			}
		case pb.CallbackType_CB_HTTP:
			config := pb.HttpCallback{}
			err = proto.Unmarshal(pcb.Config, &config)
			if err == nil {
				hcb := callback.Http(config.Url)
				if config.Headers != nil {
					for name, value := range config.Headers {
						hcb.Header(name, value)
					}
				}
				cb = hcb
			}
		case pb.CallbackType_CB_FILE:
			config := pb.FileCallback{}
			err = proto.Unmarshal(pcb.Config, &config)
			if err == nil {
				cb = callback.FileCallback(config.Path)
			}
		}
	}
	if err != nil {
		log.Printf("Parse callback error: %s", err)
		err = nil
	}
	// Parse task
	t := engine.Task(nil)
	switch params.Task.Type {
	case pb.TaskType_TASK_COMMAND:
		config := pb.CommandTask{}
		if err = proto.Unmarshal(params.Task.Config, &config); err != nil {
			break
		}
		ct := task.Command(config.Name, config.Args...)
		if config.Input != nil && len(config.Input) > 0 {
			ct.Input(config.Input)
		}
		t = ct
	case pb.TaskType_TASK_HTTP:
		config := pb.HttpTask{}
		if err = proto.Unmarshal(params.Task.Config, &config); err != nil {
			break
		}
		ht := task.Http(config.Url, config.Method, config.Body)
		if config.Headers != nil {
			for name, value := range config.Headers {
				ht.Header(name, value)
			}
		}
		t = ht
	}
	if err != nil {
		log.Printf("Parse task error: %s", err)
	} else {
		err = s.e.JobAdd(params.Queue, t, cb)
	}
	return
}
