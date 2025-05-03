//go:build !windows

package daemon

import (
	"log"
	"net"
)

func getConn() net.Conn {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	return conn
}
