package task

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os/exec"
)

var errNoCommand = errors.New("no command to execute")

type CommandTask struct {
	baseTask
	name  string
	args  []string
	input []byte
}

func (t *CommandTask) Run(ctx context.Context) (err error) {
	if t.name == "" {
		return errNoCommand
	}
	// Make command
	cmd := exec.CommandContext(ctx, t.name, t.args...)
	// Input data
	if t.input != nil {
		cmd.Stdin = bytes.NewReader(t.input)
	}
	// Capture stdout when need
	stdout := (*bytes.Buffer)(nil)
	if t.cb != nil {
		stdout = &bytes.Buffer{}
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = ioutil.Discard
	}
	cmd.Stderr = ioutil.Discard
	// Execute command
	if err = cmd.Run(); err == nil {
		if t.cb != nil && stdout != nil {
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
