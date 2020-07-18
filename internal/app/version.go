package app

import (
	"fmt"
	"runtime"
)

const (
	name    = "doppelganger"
	version = "0.0.1"
)

func Version() string {
	return fmt.Sprintf("%s %s (%s %s/%s)",
		name, version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
