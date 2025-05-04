package daemon

import (
	"log"
	"net/rpc"
)

// CallDaemonIfOpen returns the passphrase served by the RCW daemon
// (if one is available). If no RCW daemon is accessible, nil is returned.
func CallDaemonIfOpen() []byte {
	if daemonIsOpen() {
		call()
		return call()
	}
	return nil
}

// call connects to the RPC server and requests the passphrase.
func call() []byte {
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
	return []byte(reply)
}
