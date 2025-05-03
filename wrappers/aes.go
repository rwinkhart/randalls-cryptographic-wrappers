package wrappers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const (
	nonceSizeAES = 12 // GCM standard nonce size is 12 bytes
)

// EncryptAES encrypts data using AES-256-GCM.
func encryptAES(decBytes []byte, passphrase []byte) []byte {
	// generate a random salt
	salt := make([]byte, saltSize)
	io.ReadFull(rand.Reader, salt)

	// derive key from passphrase using the salt
	key := deriveKey(passphrase, salt)

	// create AES-256 cipher
	block, _ := aes.NewCipher(key)

	// create GCM mode
	aesGCM, _ := cipher.NewGCM(block)

	// generate a random nonce
	nonce := make([]byte, nonceSizeAES)
	io.ReadFull(rand.Reader, nonce)

	// encrypt the data
	ciphertext := aesGCM.Seal(nil, nonce, decBytes, nil)

	// format: salt + nonce + ciphertext
	result := make([]byte, 0, saltSize+nonceSizeAES+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result
}

// DecryptAES decrypts data using AES256-GCM.
func decryptAES(encBytes []byte, passphrase []byte) ([]byte, error) {
	if len(encBytes) < saltSize+nonceSizeAES {
		return nil, errors.New("AES256-GCM: Encrypted data is too short")
	}

	// extract salt, nonce, and ciphertext
	salt := encBytes[:saltSize]
	nonce := encBytes[saltSize : saltSize+nonceSizeAES]
	ciphertext := encBytes[saltSize+nonceSizeAES:]

	// derive key from passphrase using the salt
	key := deriveKey(passphrase, salt)

	// create AES-256 cipher
	block, _ := aes.NewCipher(key)

	// create GCM mode
	aesGCM, _ := cipher.NewGCM(block)

	// decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
