package callback

import (
	"github.com/deadblue/gostream/quietly"
	"log"
	"os"
)

type FileCallback string

func (c FileCallback) Send(result []byte) {
	file, err := os.OpenFile(string(c), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Create file error: %s", err)
		return
	}
	defer quietly.Close(file)
	_, err = file.Write(result)
	if err != nil {
		log.Printf("Write file error: %s", err)
	}
}
