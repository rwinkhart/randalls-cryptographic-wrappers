package wrappers

import (
	"runtime"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)

const (
	// parameters for Argon2
	argonTime   = 8          // set to pass 1-second test in dev environment
	argonMemory = 384 * 1024 // 384 MB

	// parameters for Scrypt
	scryptN = 262144
	scryptR = 8 // block size, recommended to be 8
	scryptP = 1 // parallelization factor, recommended to be 1

	// general constants
	keyLen   = 32 // 256 bits, key length for both algorithms
	saltSize = 16 // 128 bits, recommended salt size for both algorithms
)

var threads = uint8(runtime.NumCPU())

// DeriveKeyArgon2 derives an encryption key from a passphrase using Argon2.
func deriveKeyArgon2(passphrase []byte, salt []byte) []byte {
	return argon2.IDKey(passphrase, salt, argonTime, argonMemory, threads, keyLen)
}

// DeriveKeyScrypt derives an encryption key from a passphrase using Scrypt.
func deriveKeyScrypt(passphrase []byte, salt []byte) []byte {
	key, _ := scrypt.Key(passphrase, salt, scryptN, scryptR, scryptP, keyLen)
	return key
}
