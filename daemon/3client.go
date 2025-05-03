package daemon

import (
	"log"
	"net/rpc"
)

// Call connects to the RPC server and requests the passphrase.
func Call() string {
	// connect to the UNIX domain socket/Windows named pipe
	conn := getConn()
	defer conn.Close()

	// create an RPC client using the connection
	client := rpc.NewClient(conn)
	defer client.Close()

	// request the passphrase from the RPC server
	var reply string
	if err := client.Call("RCWService.GetPass", "hi", &reply); err != nil {
		log.Fatalf("Error calling RCWService.GetPass: %v", err)
	}

	// return the passphrase
	return reply
}
