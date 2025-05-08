package wrappers

import (
	"runtime"

	"golang.org/x/crypto/argon2"
)

const (
	// parameters for Argon2
	argonTime   = 8          // set to pass 1-second test in dev environment
	argonMemory = 384 * 1024 // 384 MB

	// general constants
	keyLen   = 32 // 256 bits, key length for both algorithms
	saltSize = 16 // 128 bits, recommended salt size for both algorithms
)

var argonThreads = uint8(runtime.NumCPU())

// DeriveKey derives an encryption key from a passphrase using Argon2.
func deriveKey(passphrase []byte, salt []byte) []byte {
	return argon2.IDKey(passphrase, salt, argonTime, argonMemory, argonThreads, keyLen)
}
