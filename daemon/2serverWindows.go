//go:build windows

package daemon

import (
	"bytes"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Microsoft/go-winio" // For Windows named pipes
	"github.com/rwinkhart/peercred-mini"
	"golang.org/x/sys/windows"
)

const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
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

	// capture sigterms to ensure listener is closed
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	// create 3-minute inactivity timer
	timer := time.NewTimer(3 * time.Minute)
	killTimer := make(chan struct{})
	go func() {
		select {
		case <-timer.C:
			log.Println("Three minutes have passed without any connections. Exiting...")
			listener.Close()
			os.Exit(0)
		case <-killTimer:
			return
		case <-sigChan:
			listener.Close()
			os.Exit(0)
		}
	}()

	// accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			close(killTimer)
			continue
		}

		// reset timer after connection is accepted
		timer.Reset(3 * time.Minute)

		// use a goroutine to check the client's identity
		go handleConn(conn)
	}
}

// handleConn verifies the identity of the client.
// It gets the PID of the client process and verifies it's running the same binary
func handleConn(conn net.Conn) {
	ucred := peercred.Get(conn)

	// get server SID (UID)
	var token windows.Token
	windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	defer token.Close()
	user, _ := token.GetTokenUser()

	// check if the RPC call is coming from an identical binary and from the same user
	callingBinPath := pidToPath(uint32(ucred.PID))
	if ucred.UID == user.User.Sid.String() && bytes.Equal(getFileHash(callingBinPath), daemonHash) {
		rpc.ServeConn(conn)
	} else {
		// invalid client; close the connection w/o a response,
		// log the client's path, and kill the daemon
		conn.Close()
		log.Printf("Request received from invalid client: PID(%d), UID(%s), Path(%s)", ucred.PID, ucred.UID, callingBinPath) // TODO log to file
		os.Exit(2)
	}
}
