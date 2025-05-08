package wrappers

// Decrypt decrypts the provided byte slice using the provided passphrase.
// Do NOT use directly if using the RCW daemon.
func Decrypt(encBytes []byte, passphrase []byte) ([]byte, error) {
	var err error = nil
	encBytes, err = decryptCha(encBytes, passphrase)
	if err != nil {
		return nil, err
	}
	encBytes, err = decryptAES(encBytes, passphrase)
	if err != nil {
		return nil, err
	}
	return encBytes, err
}

// Encrypt encrypts the provided byte slice using the provided passphrase.
// Do NOT use directly if using the RCW daemon.
func Encrypt(decBytes []byte, passphrase []byte) []byte {
	decBytes = encryptAES(decBytes, passphrase)
	decBytes = encryptCha(decBytes, passphrase)
	return decBytes
}
