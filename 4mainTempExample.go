package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "client":
			fmt.Println("RPC Reply: " + CallDaemon())
		default:
			StartDaemon()
		}
	} else {
		StartDaemon()
	}
}
