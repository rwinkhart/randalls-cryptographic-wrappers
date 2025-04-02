package main

import (
	"fmt"
	"os"
	"rcw/daemon"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "client":
			fmt.Println("RPC Reply: " + daemon.Call())
		default:
			daemon.Run()
		}
	} else {
		daemon.Run()
	}
}
