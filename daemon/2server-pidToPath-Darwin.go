//go:build darwin

package daemon

import (
	"os/exec"
	"strconv"
	"strings"
)

// pidToPath returns the path of the executable that has the given PID.
// TODO remove reliance on system command OR verify authenticity of "lsof" binary
func pidToPath(pid int) string {
	pidString := strconv.Itoa(pid)
	cmd := exec.Command("lsof", "-a", "-dtxt", "-p"+pidString)
	output, _ := cmd.Output()
	line := strings.Split(string(output), "\n")[1]
	return line[strings.Index(line, "/"):]
}
