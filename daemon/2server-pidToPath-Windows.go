//go:build windows

package daemon

import (
	"log"
	"syscall"

	"golang.org/x/sys/windows"
)

// pidToPath returns the path of the executable that has the given PID.
func pidToPath(pid uint32) string {
	// get a handle to the process
	const PROCESS_QUERY_INFORMATION = 0x0400
	const PROCESS_VM_READ = 0x0010
	hProcess, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, false, pid)
	if err != nil {
		log.Fatalf("Failed to open process %d: %v", pid, err)
	}
	defer windows.CloseHandle(hProcess)

	// query the process executable path
	var pathBuf [syscall.MAX_PATH]uint16
	pathLen := uint32(len(pathBuf))
	windows.QueryFullProcessImageName(hProcess, 0, &pathBuf[0], &pathLen)

	return syscall.UTF16ToString(pathBuf[:pathLen])
}
