package engine

import "context"

type Callback interface {
	Send(result []byte)
}

type Task interface {
	Run(ctx context.Context) error
}
