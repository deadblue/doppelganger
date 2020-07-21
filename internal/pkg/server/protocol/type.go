package protocol

type TaskType string

func (t TaskType) IsCommand() bool {
	return t == "command"
}

func (t TaskType) IsHttp() bool {
	return t == "http"
}

type CallbackType string

func (t CallbackType) IsCommand() bool {
	return t == "command"
}

func (t CallbackType) IsHttp() bool {
	return t == "http"
}
