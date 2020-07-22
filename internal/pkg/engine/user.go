// +build !windows

package engine

import (
	"log"
	"os"
)

func checkUser() {
	if os.Geteuid() == 0 {
		log.Print("It's dangerous to run doppelganger as root.")
	}
}
