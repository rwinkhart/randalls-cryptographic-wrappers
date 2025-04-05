//go:build windows

package daemon

import (
	"log"
	"net/rpc"

	"github.com/Microsoft/go-winio"
)

// Call connects to the RPC server and requests the passphrase.
func Call() string {
	// connect to the Windows named pipe
	conn, err := winio.DialPipe(socketPath, nil)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	defer conn.Close()

	// create an RPC client using the connection
	client := rpc.NewClient(conn)
	defer client.Close()

	// request the passphrase from the RPC server
	var reply string
	err = client.Call("RCWService.GetPass", "hi", &reply)
	if err != nil {
		log.Fatalf("Error calling RCWService.GetPass: %v", err)
	}

	// return the passphrase
	return reply
}
