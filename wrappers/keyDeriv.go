package wrappers

import (
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

// parameters for Argon2
const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = chacha20poly1305.KeySize
	saltSize     = 16
)

// DeriveKey derives an encryption key from a passphrase using Argon2.
func deriveKey(passphrase []byte, salt []byte) []byte {
	return argon2.IDKey(passphrase, salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}
