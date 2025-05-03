//go:build windows

package daemon

import (
	"log"
	"net"

	"github.com/Microsoft/go-winio"
)

func getConn() net.Conn {
	conn, err := winio.DialPipe(socketPath, nil)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	return conn
}
