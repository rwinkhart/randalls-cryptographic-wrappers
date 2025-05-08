package daemon

import (
	"log"
	"net/rpc"
)

// CallDaemonIfOpen uses the RCW daemon (if one is available) to
// decrypt and return data. If no RCW daemon is accessible, nil is returned.
func CallDaemonIfOpen(encBytes []byte) []byte {
	if daemonIsOpen() {
		return call(encBytes)
	}
	return nil
}

// call connects to the RPC server and requests the passphrase.
func call(encBytes []byte) []byte {
	// connect to the UNIX domain socket/Windows named pipe
	conn := getConn()
	defer conn.Close()

	// create an RPC client using the connection
	client := rpc.NewClient(conn)
	defer client.Close()

	// request the passphrase from the RPC server
	var reply []byte
	if err := client.Call("RCWService.DecryptRequest", encBytes, &reply); err != nil {
		log.Fatalf("Error calling RCWService.DecryptRequest: %v", err)
	}

	// return the passphrase
	return []byte(reply)
}
