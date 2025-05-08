package daemon

import (
	"log"
	"net"
	"net/rpc"
)

// DecryptWithDaemonIfOpen uses the RCW daemon (if one is available) to
// decrypt and return data. If no RCW daemon is accessible, nil is returned.
func DecryptWithDaemonIfOpen(encBytes []byte) []byte {
	if IsOpen() {
		return getDecFromDaemon(encBytes)
	}
	return nil
}

// EncryptWithDaemonIfOpen uses the RCW daemon (if one is available) to
// encrypt and return data. If no RCW daemon is accessible, nil is returned.
func EncryptWithDaemonIfOpen(decBytes []byte) []byte {
	if IsOpen() {
		return getEncFromDaemon(decBytes)
	}
	return nil
}

// getDecFromDaemon requests the RCW daemon to decrypt the given data.
// It returns the decrypted data.
func getDecFromDaemon(encBytes []byte) []byte {
	conn, client := connectToDaemon()
	defer conn.Close()
	defer client.Close()

	// request decBytes from the RPC server
	var decBytes []byte
	if err := client.Call("RCWService.DecryptRequest", encBytes, &decBytes); err != nil {
		log.Fatalf("Error calling RCWService.DecryptRequest: %v", err)
	}
	return decBytes
}

// getEncFromDaemon requests the RCW daemon to encrypt the given data.
// It returns the encrypted data.
func getEncFromDaemon(decBytes []byte) []byte {
	conn, client := connectToDaemon()
	defer conn.Close()
	defer client.Close()

	// request encBytes from the RPC server
	var encBytes []byte
	if err := client.Call("RCWService.EncryptRequest", decBytes, &encBytes); err != nil {
		log.Fatalf("Error calling RCWService.EncryptRequest: %v", err)
	}
	return encBytes
}

// connectToDaemon establishes a connection to the RCW daemon.
// It returns the connection and the RPC client.
// The caller is responsible for closing the connection and client.
func connectToDaemon() (net.Conn, *rpc.Client) {
	conn := getConn()
	client := rpc.NewClient(conn)
	return conn, client
}
