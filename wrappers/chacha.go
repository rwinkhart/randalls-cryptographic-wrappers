package wrappers

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	nonceSizeCha = chacha20poly1305.NonceSizeX
)

// EncryptCha encrypts data using ChaCha20-Poly1305.
func encryptCha(data []byte, passphrase []byte) []byte {
	// generate a random salt
	salt := make([]byte, saltSize)
	io.ReadFull(rand.Reader, salt)

	// derive key from passphrase using the salt
	key := deriveKey(passphrase, salt)

	// create ChaCha20-Poly1305 cipher
	stream, _ := chacha20poly1305.NewX(key)

	// generate a random nonce
	nonce := make([]byte, nonceSizeCha)
	io.ReadFull(rand.Reader, nonce)

	// encrypt the data
	ciphertext := stream.Seal(nil, nonce, data, nil)

	// format: salt + nonce + ciphertext
	result := make([]byte, 0, saltSize+nonceSizeCha+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result
}

// DecryptCha decrypts data using ChaCha20-Poly1305.
func decryptCha(encryptedData []byte, passphrase []byte) ([]byte, error) {
	if len(encryptedData) < saltSize+nonceSizeCha {
		return nil, errors.New("ChaCha20-Poly1305: Encrypted data is too short")
	}

	// extract salt, nonce, and ciphertext
	salt := encryptedData[:saltSize]
	nonce := encryptedData[saltSize : saltSize+nonceSizeCha]
	ciphertext := encryptedData[saltSize+nonceSizeCha:]

	// derive key from passphrase using the salt
	key := deriveKey(passphrase, salt)

	// create ChaCha20-Poly1305 cipher
	stream, _ := chacha20poly1305.NewX(key)

	// decrypt the data
	plaintext, err := stream.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
