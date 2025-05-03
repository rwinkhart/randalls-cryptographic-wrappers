package wrappers

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	nonceSizeCha = chacha20poly1305.NonceSizeX
	saltSizeCha  = 16
)

// EncryptCha encrypts data using ChaCha20-Poly1305
func EncryptCha(data []byte, passphrase []byte) []byte {
	// generate a random salt
	salt := make([]byte, saltSizeCha)
	io.ReadFull(rand.Reader, salt)

	// derive key from passphrase using the salt
	// TODO ensure the passphrase is consistent (store a hashed version to compare against)
	key := deriveKey(passphrase, salt)

	// create ChaCha20-Poly1305 cipher
	aead, _ := chacha20poly1305.NewX(key)

	// generate a random nonce
	nonce := make([]byte, nonceSizeCha)
	io.ReadFull(rand.Reader, nonce)

	// encrypt the data
	ciphertext := aead.Seal(nil, nonce, data, nil)

	// format: salt + nonce + ciphertext
	result := make([]byte, 0, saltSizeCha+nonceSizeCha+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result
}

// DecryptCha decrypts data using ChaCha20-Poly1305
func DecryptCha(encryptedData []byte, passphrase []byte) []byte {
	if len(encryptedData) < saltSizeCha+nonceSizeCha {
		fmt.Println("Encrypted data is too short")
		return nil
	}

	// extract salt, nonce, and ciphertext
	salt := encryptedData[:saltSizeCha]
	nonce := encryptedData[saltSizeCha : saltSizeCha+nonceSizeCha]
	ciphertext := encryptedData[saltSizeCha+nonceSizeCha:]

	// derive key from passphrase using the salt
	key := deriveKey(passphrase, salt)

	// create ChaCha20-Poly1305 cipher
	aead, _ := chacha20poly1305.NewX(key)

	// decrypt the data
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Printf("Decryption failed (possibly wrong passphrase): %s", err.Error())
		return nil
	}

	return plaintext
}
