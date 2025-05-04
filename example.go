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
// rcw init <passwd> : Generates the required sanity check file
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

const (
	outputFile = "ex-cipher.rcw"
	sanityFile = "ex-sanity.rcw"
)

func main() {
	switch len(os.Args) {
	case 2:
		// serve data
		daemon.Start(os.Args[1])
	case 3:
		if os.Args[1] == "init" {
			// create sanity check file
			err := wrappers.GenSanityCheck(sanityFile, []byte(os.Args[2]))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// decrypt file
		encBytes, _ := os.ReadFile(outputFile)
		decBytes, err := wrappers.Decrypt(encBytes, []byte(os.Args[2]))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(decBytes))
	case 4:
		// encrypt data (from cli args)
		err := wrappers.RunSanityCheck(sanityFile, []byte(os.Args[3]))
		if err != nil {
			fmt.Println(err)
			return
		}
		encBytes := wrappers.Encrypt([]byte(os.Args[2]), []byte(os.Args[3]))
		os.WriteFile(outputFile, encBytes, 0600)
	default:
		// request served data
		fmt.Println(daemon.Call())
	}
}
