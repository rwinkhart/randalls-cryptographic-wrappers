package daemon

import (
	"os"
)

var binPath, _ = os.Executable() // store binary path

// daemonIsOpen checks if the socket/named pipe for the rcw
// daemon exists and returns a boolean indicator.
func daemonIsOpen() bool {
	fileInfo, err := os.Stat(socketPath)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}
