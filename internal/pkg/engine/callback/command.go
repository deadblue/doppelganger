package callback

import (
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"io"
	"io/ioutil"
	"os/exec"
)

type CommandCallback struct {
	name string
	args []string
}

func (cc *CommandCallback) Send(r io.Reader) error {
	cmd := exec.Command(cc.name, cc.args...)
	cmd.Stdin = r
	cmd.Stdout, cmd.Stderr = ioutil.Discard, ioutil.Discard
	return cmd.Run()
}

func Command(name string, args ...string) engine.Callback {
	return &CommandCallback{
		name: name,
		args: append([]string{}, args...),
	}
}
