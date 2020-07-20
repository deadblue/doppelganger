package task

import (
	"bytes"
	"context"
	"io/ioutil"
	"os/exec"
)

type CommandTask struct {
	baseTask
	name  string
	args  []string
	input []byte
}

func (t *CommandTask) Run(ctx context.Context) (err error) {
	// Make request
	cmd := exec.CommandContext(ctx, t.name, t.args...)
	// Input data
	if t.input != nil {
		cmd.Stdin = bytes.NewReader(t.input)
	}
	// Capture stdout when need
	if t.cb != nil {
		cmd.Stdout = &bytes.Buffer{}
	} else {
		cmd.Stdout = ioutil.Discard
	}
	cmd.Stderr = ioutil.Discard
	// Run command
	if err = cmd.Run(); err == nil {
		if t.cb != nil {
			stdout := cmd.Stdout.(*bytes.Buffer)
			go t.cb.Send(stdout.Bytes())
		}
	}
	return nil
}

func (t *CommandTask) Input(data []byte) *CommandTask {
	t.input = make([]byte, len(data))
	copy(t.input, data)
	return t
}

// Command create a command based task.
func Command(name string, args []string) *CommandTask {
	if args == nil {
		args = []string{}
	}
	return &CommandTask{
		name: name,
		args: args,
	}
}
