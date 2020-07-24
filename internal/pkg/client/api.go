package client

import (
	"github.com/deadblue/doppelganger/internal/pkg/message"
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func (c *Client) AddCommandTask(queue string, cmdName string, cmdArgs ...string) (err error) {
	config, err := proto.Marshal(&pb.CommandTask{
		Name: cmdName,
		Args: cmdArgs,
	})
	if err != nil {
		return
	}
	return c.call(message.MethodAddTask, &pb.AddTaskParams{
		Queue:  queue,
		Type:   pb.TaskType_COMMAND_TASK,
		Config: config,
	})
}

func (c *Client) AddHttpTask(queue string, url string) (err error) {
	config, err := proto.Marshal(&pb.HttpTask{
		Url:    url,
		Method: http.MethodGet,
	})
	if err != nil {
		return
	}
	return c.call(message.MethodAddTask, &pb.AddTaskParams{
		Queue:  queue,
		Type:   pb.TaskType_HTTP_TASK,
		Config: config,
	})
}
