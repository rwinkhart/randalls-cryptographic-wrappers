package main

import (
	"fmt"
	"os"
	"rcw/daemon"
	"rcw/wrappers"
)

// This sample program serves purley as a way to interactively test the features
// of RCW before building it into your own application.
//
// Usage:
// rcw <text> : Runs the rcw daemon to serve the provided text for three minutes
// rcw : Requests the data served by the RCW daemon and outputs it to stdout
// rcw enc <text> <passwd> : Encrypts the provided text and outputs the ciphertext to encrypted-example.txt
// rcw dec <passwd> : Decrypts the provided file and outputs the plaintext to stdout

func main() {
	switch len(os.Args) {
	case 2:
		// serve data
		daemon.Start(os.Args[1])
	case 3:
		// decrypt file
		encBytes, _ := os.ReadFile("encrypted-example.txt")
		decBytes := wrappers.DecryptCha(encBytes, []byte(os.Args[2]))
		fmt.Println(string(decBytes))
	case 4:
		// encrypt data (from cli args)
		encBytes := wrappers.EncryptCha([]byte(os.Args[2]), []byte(os.Args[3]))
		os.WriteFile("encrypted-example.txt", encBytes, 0644)
	default:
		// request served data
		fmt.Println(daemon.Call())
	}
}
