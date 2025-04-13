package daemon

import (
	"crypto/sha256"
	"io"
	"os"
)

var daemonHash []byte
var passphrase string

// RCWService provides an RPC method.
type RCWService struct{}

// GetPass is the RPC method.
// For now (as a test/example), it returns "hello" if the input is "hi".
func (h *RCWService) GetPass(request string, reply *string) error {
	*reply = passphrase
	return nil
}

// getFileHash returns the SHA256 hash of the file at the given path.
func getFileHash(path string) []byte {
	file, _ := os.Open(path)
	hash := sha256.New()
	io.Copy(hash, file)
	return hash.Sum(nil)
}
