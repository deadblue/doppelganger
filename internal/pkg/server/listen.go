package server

import (
	"net"
	"os"
)

func Listen(network, address string) (l net.Listener, err error) {
	if network == "unix" || network == "unixpacket" {
		// Try remove socket file before listen
		if err = os.Remove(address); !os.IsNotExist(err) {
			return
		}
	}
	l, err = net.Listen(network, address)
	if err == nil {
		if ul, ok := l.(*net.UnixListener); ok {
			ul.SetUnlinkOnClose(true)
		}
	}
	return
}
