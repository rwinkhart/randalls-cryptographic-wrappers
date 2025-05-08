package daemon

import (
	"log"
	"net"
	"net/rpc"
)

// GetDec requests the RCW daemon to decrypt the given data.
// It returns the decrypted data.
func GetDec(encBytes []byte) []byte {
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

// GetEnc requests the RCW daemon to encrypt the given data.
// It returns the encrypted data.
func GetEnc(decBytes []byte) []byte {
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
