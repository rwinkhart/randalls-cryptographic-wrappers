//go:build android && termux && !interactive

package daemon

import "path/filepath"

var socketPath = "/data/data/com.termux/files/usr/tmp/" + filepath.Base(binPath) + "-rcwd.sock" // store UNIX socket path
