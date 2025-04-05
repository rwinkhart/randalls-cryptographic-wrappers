//go:build windows

package daemon

import "path/filepath"

var socketPath = `\\.\pipe\` + filepath.Base(binPath) + `-rcwd` // store Windows named pipe path
