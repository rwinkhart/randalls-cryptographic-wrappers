//go:build windows

package daemon

import (
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/Microsoft/go-winio" // For Windows named pipes
)

const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
)

// Run should be called to start an RPC server using Windows named pipes
func Run() {
	// register RCWService with the RPC package
	if err := rpc.Register(&RCWService{}); err != nil {
		log.Fatalf("Error registering RPC service: %v", err)
	}

	// store the hash of the daemon binary
	daemonHash = getFileHash(binPath)

	// configure the named pipe
	pipeConfig := &winio.PipeConfig{
		SecurityDescriptor: "", // Default security
		MessageMode:        true,
		InputBufferSize:    65536,
		OutputBufferSize:   65536,
	}

	// create the named pipe listener
	listener, err := winio.ListenPipe(socketPath, pipeConfig)
	if err != nil {
		log.Fatalf("Failed to listen on named pipe %s: %v", socketPath, err)
	}
	defer listener.Close()
	log.Printf("RPC daemon listening on %s", socketPath)

	// accept connections (timeout after 3 minutes of inactivity)
	for {
		// set deadline for accepting new connections
		//listener.SetDeadline(time.Now().Add(3 * time.Minute)) TODO FIX

		conn, err := listener.Accept()
		if err != nil {
			if os.IsTimeout(err) {
				log.Println("Three minutes have passed without any connections. Exiting...")
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
// It gets the PID of the client process and verifies it's running the same binary
// TODO ADD AUTHENTICATION
func handleConn(conn net.Conn) {
	rpc.ServeConn(conn)
}
