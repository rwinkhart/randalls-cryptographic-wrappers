package daemon

import "os"

func daemonIsOpen() bool {
	// check if socketPath exists and is a file
	fileInfo, err := os.Stat(socketPath)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}
