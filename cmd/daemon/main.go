package main

import (
	"github.com/deadblue/doppelganger/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Panicf("Exit with error: %s", err)
	}
}
