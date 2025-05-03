package wrappers

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

func Encrypt(decBytes []byte, passphrase []byte) []byte {
	decBytes = encryptAES(decBytes, passphrase)
	decBytes = encryptCha(decBytes, passphrase)
	return decBytes
}
