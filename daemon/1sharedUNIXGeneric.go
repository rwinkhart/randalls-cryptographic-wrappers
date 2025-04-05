//go:build !windows

package daemon

import (
	"path/filepath"
)

var socketPath = "/tmp/" + filepath.Base(binPath) + "-rcwd.sock" // store UNIX socket path
