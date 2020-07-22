// Package suid provides function to generate a Sortable-Unquie-IDentifier.
package suid

import (
	"hash/crc32"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	charBase32 = "0123456789abcdefghijkmnpqrtuvwxy"
)

var (
	machineId []byte
	processId []byte

	locker = &sync.Mutex{}

	currPrefix = ""
	currSeq    = 0
)

func init() {
	// Machine ID
	h := crc32.NewIEEE()
	if nifs, err := net.Interfaces(); err == nil {
		for _, nif := range nifs {
			_, _ = h.Write(nif.HardwareAddr)
		}
	}
	machineId = h.Sum(nil)
	// Process ID
	pid := os.Getpid()
	processId = []byte{
		byte(pid >> 16), byte(pid >> 8), byte(pid),
	}
}

func getSequence(prefix string) int {
	locker.Lock()
	defer locker.Unlock()

	if prefix == currPrefix {
		currSeq += 1
	} else {
		currPrefix = prefix
		currSeq = 0
	}
	return currSeq
}

func encode(buf []byte) string {
	sb := new(strings.Builder)
	for i := 0; i < 3; i++ {
		n := big.NewInt(0).
			SetBytes(buf[i*5 : (i+1)*5]).
			Uint64()
		for j := 0; j < 8; j++ {
			index := (n >> ((7 - j) * 5)) & 0x1f
			sb.WriteByte(charBase32[index])
		}
	}
	return sb.String()
}

func NextRaw() []byte {
	// timestamp in milliseconds
	now := time.Now().UnixNano() / 1000000
	timestamp := []byte{
		byte(now >> 40), byte(now >> 32), byte(now >> 24),
		byte(now >> 16), byte(now >> 8), byte(now),
	}
	// ID prefix
	buf := make([]byte, 15)
	copy(buf, timestamp)
	copy(buf[6:], machineId)
	copy(buf[10:], processId)
	// get sequence
	index := getSequence(string(buf))
	buf[13], buf[14] = byte(index>>8), byte(index)

	return buf
}

func Next() string {
	return encode(NextRaw())
}
