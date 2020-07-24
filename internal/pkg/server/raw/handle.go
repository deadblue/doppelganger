package raw

import (
	"errors"
	"github.com/deadblue/doppelganger/internal/pkg/message"
	"github.com/deadblue/doppelganger/internal/pkg/message/pb"
	"google.golang.org/protobuf/proto"
	"log"
)

var (
	errUnknownMethod = errors.New("unknown method")
)

func (s *Server) handle(req *pb.Request) (result proto.Message, err error) {
	switch req.Method {
	case message.MethodAddTask:
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
	log.Printf("AddTask params: %s", params)
	return
}
