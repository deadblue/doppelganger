package callback

import (
	"github.com/deadblue/doppelganger/internal/pkg/engine"
	"github.com/deadblue/gostream/quietly"
	"io"
	"os"
)

type FileCallback string

func (fc FileCallback) Send(r io.Reader) (err error) {
	file, err := os.OpenFile(string(fc), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer quietly.Close(file)
	_, err = io.Copy(file, r)
	return
}

func File(path string) engine.Callback {
	return FileCallback(path)
}
