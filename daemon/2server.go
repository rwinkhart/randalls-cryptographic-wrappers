package daemon

import (
	"crypto/sha256"
	"io"
	"os"

	"github.com/rwinkhart/rcw/wrappers"
)

var daemonHash []byte
var globalPassphrase []byte

// RCWService provides an RPC method.
type RCWService struct{}

// DecryptRequest is the RPC method that decrypts the incoming data using
// the global passphrase and returns the decrypted data
func (h *RCWService) DecryptRequest(request []byte, reply *[]byte) error {
	var err error
	*reply, err = wrappers.Decrypt(request, globalPassphrase)
	if err != nil {
		return err
	}
	return nil
}

// getFileHash returns the SHA256 hash of the file at the given path.
func getFileHash(path string) []byte {
	file, _ := os.Open(path)
	hash := sha256.New()
	io.Copy(hash, file)
	return hash.Sum(nil)
}
