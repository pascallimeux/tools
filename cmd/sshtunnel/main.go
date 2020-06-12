package main

import (
	"log"
	"os"
	"time"
	"tools/pkg/sshtunnel"

	"golang.org/x/crypto/ssh"
)

func main() {
	tunnel := sshtunnel.NewSSHTunnel(
		// User and host of tunnel server, it will default to port 22
		// if not specified.
		"pascal@192.168.20.120",

		// Pick ONE of the following authentication methods:
		//sshtunnel.PrivateKeyFile("path/to/private/key.pem"), // 1. private key
		ssh.Password("pascal"), // 2. password

		// The destination host and port of the actual server.
		"192.168.20.53:8080",

		// The local port you want to bind the remote port to.
		// Specifying "0" will lead to a random port.
		"8888",
	)

	// You can provide a logger for debugging, or remove this line to
	// make it silent.
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	// Start the server in the background. You will need to wait a
	// small amount of time for it to bind to the localhost port
	// before you can start sending connections.
	go tunnel.Start()
	time.Sleep(10000 * time.Millisecond)

}
