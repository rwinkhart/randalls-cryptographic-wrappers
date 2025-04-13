package main

import (
	"fmt"
	"os"
	"rcw/daemon"
)

func main() {
	if len(os.Args) > 1 {
		daemon.Start(os.Args[1])
	} else {
		fmt.Println(daemon.Call())
	}
}
