package main

import (
	"log"
	"net"
	"net/rpc"
)

func CallDaemon() string {
	// connect to the UNIX domain socket
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	defer conn.Close()

	// create an RPC client using the connection
	client := rpc.NewClient(conn)
	defer client.Close()

	// request the passphrase from the RPC server
	var reply string
	err = client.Call("RCWService.GetPass", "hi", &reply)
	if err != nil {
		log.Fatalf("Error calling RCWService.GetPass: %v", err)
	}

	// return the passphrase
	return reply
}
