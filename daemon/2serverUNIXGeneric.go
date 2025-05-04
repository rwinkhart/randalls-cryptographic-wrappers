//go:build !windows

package daemon

import (
	"bytes"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	peercred "github.com/rwinkhart/peercred-mini"
)

// Start should be called to serve the given passphrase through an RPC daemon.
func Start(passphrase string) {
	// store passphrase to be referenced by GetPass method
	globalPassphrase = passphrase

	// register RCWService with the RPC package
	if err := rpc.Register(&RCWService{}); err != nil {
		log.Fatalf("Error registering RPC service: %v", err)
	}

	// store the hash of the daemon binary
	daemonHash = getFileHash(binPath)

	// listen on the Unix domain socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on UNIX socket %s: %v", socketPath, err)
	}
	defer listener.Close()
	log.Printf("RPC daemon listening on unix://%s", socketPath)

	// accept connections (timeout after 3 minutes of inactivity)
	for {
		listener.(*net.UnixListener).SetDeadline(time.Now().Add(3 * time.Minute))

		conn, err := listener.Accept()
		if err != nil {
			if err.(net.Error).Timeout() {
				log.Println("Three minutes have passed without any connections. Exiting...")
				listener.Close()
				os.Exit(0)
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		// use a goroutine to check the client's identity
		go handleConn(conn)
	}
}

// handleConn verifies the identity of the client.
// It uses the file descriptor of the connection to get the PID of the client,
// which is then used to get the path of the client's executable and calculate its hash.
// The passphrase is only returned if the client's executable hash matches the daemon's hash
// and if the request is coming from the same user.
// This ensures that only the binary the daemon is embedded in can retrieve the passphrase.
func handleConn(conn net.Conn) {
	ucred := peercred.Get(conn)

	// check if the RPC call is coming from an identical binary and from the same user
	callingBinPath := pidToPath(ucred.PID)
	if ucred.UID == strconv.Itoa(os.Getuid()) && bytes.Equal(getFileHash(callingBinPath), daemonHash) {
		// valid client; hand off the connection to the RPC server
		rpc.ServeConn(conn)
	} else {
		// invalid client; close the connection w/o a response,
		// log the client's path, and kill the daemon
		conn.Close()
		log.Printf("Request received from invalid client: PID(%d), UID(%s), Path(%s)", ucred.PID, ucred.UID, callingBinPath) // TODO log to file
		os.Exit(2)
	}
}
