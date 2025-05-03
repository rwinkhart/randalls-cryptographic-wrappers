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
// rcw dec <passwd> : Decrypts encrypted-example.txt and outputs the plaintext to stdout

// TODO Tests:
// Salt (aes+chacha)
// Nonce (aes+chacha)
// Encryption (individual+combined)
// Decryption (individual+combined)
// RPC password sharing

// TODO Enhancements:
// Security:
// 	   Play with nonce sizes and Argon2 parameters to find the best speed-security balance
// Standalone cmd:
//     Usable as symmetric-only GPG replacement

func main() {
	switch len(os.Args) {
	case 2:
		// serve data
		daemon.Start(os.Args[1])
	case 3:
		// decrypt file
		encBytes, _ := os.ReadFile("encrypted-example.txt")
		decBytes, err := wrappers.Decrypt(encBytes, []byte(os.Args[2]))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(decBytes))
	case 4:
		// encrypt data (from cli args)
		encBytes := wrappers.Encrypt([]byte(os.Args[2]), []byte(os.Args[3]))
		os.WriteFile("encrypted-example.txt", encBytes, 0644)
	default:
		// request served data
		fmt.Println(daemon.Call())
	}
}
