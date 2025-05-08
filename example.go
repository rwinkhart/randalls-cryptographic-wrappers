package main

import (
	"fmt"
	"os"

	"github.com/rwinkhart/rcw/daemon"
	"github.com/rwinkhart/rcw/wrappers"
	"golang.org/x/term"
)

// This sample program serves purley as a way to interactively test the features
// of RCW before building it into your own application.
//
// Usage:
// rcw init <passwd> : Generates the required sanity check file
// rcw <passphrase> : Runs the rcw daemon to decrypt data for three minutes
// rcw enc <text> <passwd> : Encrypts the provided text and outputs the ciphertext to encrypted-example.txt
// rcw dec : Decrypts ex-cipher.rcw and outputs the plaintext to stdout (attempts to use daemon, falls back to user input for passphrase)

// TODO Tests:
// Salt (aes+chacha)
// Nonce (aes+chacha)
// Encryption (individual+combined)
// Decryption (individual+combined)
// RPC password sharing

// TODO Enhancements:
// Standalone cmd:
//     Usable as symmetric-only GPG replacement

const (
	outputFile = "ex-cipher.rcw"
	sanityFile = "ex-sanity.rcw"
)

func main() {
	switch len(os.Args) {
	case 2:
		if os.Args[1] == "dec" {
			// decrypt file (using daemon if available)
			// rcw dec
			encBytes, err := os.ReadFile(outputFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			decBytes := daemon.CallDaemonIfOpen(encBytes)
			if decBytes == nil {
				fmt.Println("No RCW daemon available")
				passphrase := inputHidden("Enter RCW passphrase:")
				decBytes, err = wrappers.Decrypt(encBytes, passphrase)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			fmt.Println(string(decBytes))
			return
		}
		// run decrypter daemon
		// rcw <passwd>
		daemon.Start([]byte(os.Args[1]))
	case 3:
		if os.Args[1] == "init" {
			// create sanity check file
			// rcw init <passwd>
			err := wrappers.GenSanityCheck(sanityFile, []byte(os.Args[2]))
			if err != nil {
				fmt.Println(err)
			}
		}
	case 4:
		// encrypt data (from cli args)
		// rcw enc <text> <passwd>
		err := wrappers.RunSanityCheck(sanityFile, []byte(os.Args[3]))
		if err != nil {
			fmt.Println(err)
			return
		}
		encBytes := wrappers.Encrypt([]byte(os.Args[2]), []byte(os.Args[3]))
		os.WriteFile(outputFile, encBytes, 0600)
	default:
		fmt.Println("Usage: rcw [init <passwd>] | [enc <text> <passwd>] | dec | <passwd>")
	}
}

// inputHidden prompts the user for input and returns the input as a byte array, hiding the input from the terminal.
func inputHidden(prompt string) []byte {
	fmt.Print("\n" + prompt + " ")
	byteInput, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return byteInput
}
