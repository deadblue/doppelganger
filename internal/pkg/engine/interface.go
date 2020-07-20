package engine

import "context"

type Callback interface {
	Send(result []byte)
}

type Task interface {
	Callback(cb Callback)
	Run(ctx context.Context) error
}
