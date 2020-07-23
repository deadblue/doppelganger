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

	// For generate seq
	mu = &sync.Mutex{}

	currTs  = int64(0)
	currSeq = 0
)

type SUID [15]byte

func (id SUID) String() string {
	sb := new(strings.Builder)
	for i := 0; i < 3; i++ {
		n := big.NewInt(0).
			SetBytes(id[i*5 : (i+1)*5]).
			Uint64()
		for j := 0; j < 8; j++ {
			index := (n >> ((7 - j) * 5)) & 0x1f
			sb.WriteByte(charBase32[index])
		}
	}
	return sb.String()
}

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

func getSequence(ts int64) int {
	mu.Lock()
	defer mu.Unlock()

	if ts == currTs {
		currSeq += 1
	} else {
		currTs, currSeq = ts, 0
	}
	return currSeq
}

func Next() SUID {
	// ID prefix
	id := SUID{}
	// timestamp in milliseconds
	now := time.Now().UnixNano() / 1000000
	timestamp := []byte{
		byte(now >> 40), byte(now >> 32), byte(now >> 24),
		byte(now >> 16), byte(now >> 8), byte(now),
	}
	copy(id[0:], timestamp)
	copy(id[6:], machineId)
	copy(id[10:], processId)
	// get sequence
	seq := getSequence(now)
	id[13], id[14] = byte(seq>>8), byte(seq)

	return id
}
