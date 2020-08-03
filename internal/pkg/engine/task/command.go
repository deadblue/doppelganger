package task

import (
	"bytes"
	"context"
	"errors"
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"io"
	"io/ioutil"
	"os/exec"
)

var errNoCommand = errors.New("no command to execute")

type CommandTask struct {
	name  string
	args  []string
	input []byte
}

func (t *CommandTask) Run(ctx context.Context, cb engine.Callback) (err error) {
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
	if cb != nil {
		cmd.Stdout = &bytes.Buffer{}
	} else {
		cmd.Stdout = ioutil.Discard
	}
	cmd.Stderr = ioutil.Discard
	// Execute command
	if err = cmd.Run(); err == nil {
		// Send result to callback
		if cb != nil {
			if stdout, ok := cmd.Stdout.(io.Reader); ok {
				err = cb.Send(stdout)
			} else {
				// What the hell!?
			}
		}
	}
	return
}

func (t *CommandTask) Input(data []byte) *CommandTask {
	t.input = make([]byte, len(data))
	copy(t.input, data)
	return t
}

// Command creates a command based task.
func Command(name string, args ...string) *CommandTask {
	ct := &CommandTask{
		name: name,
		args: append([]string{}, args...),
	}
	return ct
}
