package daemon

import (
	"github.com/deadblue/gostream/quietly"
	"gopkg.in/yaml.v2"
	"os"
)

type QueueConf struct {
	CoreSize int `yaml:"core_size"`
}

type NamedQueueConf struct {
	Name   string     `yaml:"name"`
	Config *QueueConf `yaml:"config"`
}

type ListenerConf struct {
	Type string `yaml:"type"`
	Addr string `yaml:"addr"`
}

type ServerConf struct {
	Http *ListenerConf `yaml:"http"`
	Raw  *ListenerConf `yaml:"raw"`
}

type Conf struct {
	DefaultQueue *QueueConf        `yaml:"default_queue"`
	Queues       []*NamedQueueConf `yaml:"queues"`
	Server       *ServerConf       `yaml:"server"`
}

func LoadConf(path string, conf *Conf) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer quietly.Close(file)
	return yaml.NewDecoder(file).Decode(conf)
}
