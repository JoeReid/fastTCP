package main

import (
	"github.com/JoeReid/fastTCP"
	"io"
	"log"
)

// handler implements a basic echo server using the standard lib io.copy
// function.
func handler(tcpFile io.ReadWriter) {
	io.Copy(tcpFile, tcpFile)
}

func main() {
	// Set the package level logger
	fastTCP.Logger = &log.Logger{}

	// create a IPv6 TCP echo server that runs on localhost with
	// TCP_DEFER_ACCEPT enabled.
	server := fastTCP.NewServer("127.0.0.1:6543", handler, fastTCP.TCPOptions{
		DeferAccept: true,
		IPv6:        true,
	})

	// Start the TCP server
	err := server.ListenTCP()
	if err != nil {
		panic(err)
	}
}
