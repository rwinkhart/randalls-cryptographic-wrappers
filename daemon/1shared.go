package daemon

import (
	"os"
	"path/filepath"
)

var binPath, _ = os.Executable()                                 // store binary path
var socketPath = "/tmp/" + filepath.Base(binPath) + "-rcwd.sock" // store UNIX socket path
