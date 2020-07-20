package callback

import (
	"bytes"
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"log"
	"os/exec"
)

type CommandCallback struct {
	name string
	args []string
}

func (c *CommandCallback) Send(result []byte) {
	cmd := exec.Command(c.name, c.args...)
	cmd.Stdin = bytes.NewReader(result)
	if err := cmd.Run(); err != nil {
		log.Printf("Send result to callback error: %s", err)
	}
}

// Command creates a command based callback.
func Command(name string, args []string) engine.Callback {
	return &CommandCallback{
		name: name,
		args: args,
	}
}
