package wrappers

import (
	"errors"
	"os"
)

// GenSanityCheck creates an encrypted file containing known plaintext
// to later be used for ensuring the user does not encrypt data with
// an incorrect passphrase.
func GenSanityCheck(path string, passphrase []byte) error {
	err := os.WriteFile(path, Encrypt([]byte("thx4usin'rcw"), passphrase), 0600)
	return err
}

// RunSanityCheck should be run before any encryption operation
// to ensure the user does not encrypt data with an incorrect passphrase.
// Failure to perform this check could result in data loss.
func RunSanityCheck(path string, passphrase []byte) error {
	encBytes, err := os.ReadFile(path)
	if err != nil {
		return errors.New("Failed to read sanity check file (" + path + ")")
	}
	decBytes, _ := Decrypt(encBytes, passphrase)
	if string(decBytes) == "thx4usin'rcw" {
		return nil
	}
	return errors.New("Sanity check failed (likely due to inconsistent passphrase)")
}
