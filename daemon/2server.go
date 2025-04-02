package daemon

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"syscall"
	"time"
)

var daemonHash []byte

// RCWService provides an RPC method.
type RCWService struct{}

// Run should be called to start an RPC server.
func Run() {
	// store the hash of the daemon binary
	daemonHash = getFileHash(binPath)

	// remove the socket file if it already exists
	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			log.Fatalf("Failed to remove existing socket: %v", err)
		}
	}

	// register RCWService with the RPC package
	if err := rpc.Register(&RCWService{}); err != nil {
		log.Fatalf("Error registering RPC service: %v", err)
	}

	// listen on the Unix domain socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on UNIX socket %s: %v", socketPath, err)
	}
	defer listener.Close()
	log.Printf("RPC daemon listening on unix://%s", socketPath)

	// Accept connections (timeout after 3 minutes of inactivity)
	for {
		listener.(*net.UnixListener).SetDeadline(time.Now().Add(3 * time.Minute))

		conn, err := listener.Accept()
		if err != nil {
			if err.(net.Error).Timeout() {
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

// GetPass is the RPC method.
// For now (as a test/example), it returns "hello" if the input is "hi".
func (h *RCWService) GetPass(request string, reply *string) error {
	if request == "hi" {
		*reply = "hello"
		return nil
	}
	return errors.New("unexpected input, expected \"hi\"")
}

// handleConn verifies the identity of the client.
// It uses the file descriptor of the connection to get the PID of the client,
// which is then used to get the path of the client's executable and calculate its hash.
// The passphrase is only returned if the client's executable hash matches the daemon's hash
// and if the request is coming from the same user.
// This ensures that only the binary the daemon is embedded in can retrieve the passphrase.
func handleConn(conn net.Conn) {
	// ensure access to the underlying file descriptor
	unixConn, ok := conn.(*net.UnixConn)
	if !ok {
		log.Printf("Connection is not a Unix domain socket")
		conn.Close()
		return
	}
	// obtain SyscallConn for direct access to the file descriptor
	rawConn, err := unixConn.SyscallConn()
	if err != nil {
		log.Printf("SyscallConn error: %v", err)
		conn.Close()
		return
	}

	var callingPID int32
	var callingUID int
	_ = rawConn.Control(func(fd uintptr) {
		// use syscall.GetsockoptUcred to fetch credentials
		ucred, err := syscall.GetsockoptUcred(int(fd), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
		if err != nil {
			log.Printf("Error getting peer credentials: %v", err)
			return
		}
		callingPID = ucred.Pid
		callingUID = int(ucred.Uid)
	})

	// check if the RPC call is coming from an identical binary and from the same user
	callingBinPath := pidToPath(callingPID)
	if callingUID == os.Getuid() && bytes.Equal(getFileHash(callingBinPath), daemonHash) {
		// valid client; hand off the connection to the RPC server
		rpc.ServeConn(conn)
	} else {
		// invalid client; close the connection w/o a response,
		// log the client's path, and kill the daemon
		conn.Close()
		log.Printf("Request received from invalid client: PID(%d), UID(%d), Path(%s)", callingPID, callingUID, callingBinPath) // TODO log to file
		os.Exit(2)
	}
}

// pidToPath returns the path of the executable that has the given PID.
func pidToPath(pid int32) string {
	path, _ := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	return path
}

// getFileHash returns the SHA256 hash of the file at the given path.
func getFileHash(path string) []byte {
	file, _ := os.Open(path)
	hash := sha256.New()
	io.Copy(hash, file)
	return hash.Sum(nil)
}
