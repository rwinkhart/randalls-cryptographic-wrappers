package wrappers

import (
	"runtime"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// parameters for Argon2
	argonTime   = 1
	argonMemory = 64 * 1024
	argonKeyLen = chacha20poly1305.KeySize

	// general constants
	saltSize = 16
)

var argonThreads = uint8(runtime.NumCPU())

// DeriveKey derives an encryption key from a passphrase using Argon2.
func deriveKey(passphrase []byte, salt []byte) []byte {
	return argon2.IDKey(passphrase, salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}
