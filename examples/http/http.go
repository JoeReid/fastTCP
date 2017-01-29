package main

import (
	"fmt"
	"github.com/JoeReid/fastTCP"
	fastHTTP "github.com/JoeReid/fastTCP/http"
	"net/http"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HTTP test")
}

func main() {
	httpServer := fastHTTP.NewHTTPServer(
		http.HandlerFunc(httpHandler),
	)

	tcpServer := fastTCP.NewServer(
		":6543",
		httpServer.NewConn,
		fastTCP.TCPOptions{},
	)

	err := tcpServer.ListenTCP()
	if err != nil {
		panic(err)
	}
}
