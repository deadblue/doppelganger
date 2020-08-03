package engine

import (
	"context"
	"io"
)

type Callback interface {
	Send(r io.Reader) (err error)
}

type Task interface {
	Run(ctx context.Context, cb Callback) error
}
