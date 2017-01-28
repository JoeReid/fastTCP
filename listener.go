package fastTCP

import (
	"github.com/valyala/tcplisten"
	"io"
	"net"
	"runtime"
	"strings"
	"time"
)

// Printer is a simple one method interface to allow a caller to provide a
// custom logger to this package.
type Printer interface {
	Printf(format string, v ...interface{})
}

// Logger is a custom logger to use for this package. Logger is nil by default.
var Logger Printer

// TCPOptions is used to configure optional performance tweaks to the TCP
// socket listener.
//
// DeferAccept corresponds to the TCP_DEFER_ACCEPT flag. If true the listener
// will set up the socket for this behaviour. See
// http://man7.org/linux/man-pages/man7/tcp.7.html for details
//
// FastOpen corresponds to the TCP_FASTOPEN flag. If true the listener will set
// up the socket for this behaviour. See https://lwn.net/Articles/508865/ for
// details.
type TCPOptions struct {
	DeferAccept bool
	FastOpen    bool
}

type Server struct {
	laddr   string
	handler func(io.ReadWriter)
	options TCPOptions
}

// NewServer returns a new Server instance configured to serve on a given local
// address with performance tweaks as defined in the provided TCPOptions.
//
// The Server will call the provided handler function for each new TCP
// connection passing the underlying os.File as an argument.
//
// The syntax of laddr is "host:port", like "127.0.0.1:8080".
// If host is omitted, as in ":8080", all available interfaces are used instead
// of just the interface with the given host address.
//
// The Server will attempt to open multiple tcp listeners (one per CPU core)
// using the SO_REUSEPORT socket option if the OS supports it. See
// http://man7.org/linux/man-pages/man7/socket.7.html for details.
//
// If the OS doesn't support this socket option, or fails to setup the socket
// like this for some other reason, the Server will degrade to using the
// standard library net.Listener implementation.
func NewServer(laddr string, handler func(io.ReadWriter), options TCPOptions) *Server {
	return &Server{
		laddr:   laddr,
		handler: handler,
		options: options,
	}
}

// spawnListener returns a new net.Listener built with the performance tweaks
// defined in t.options.
func (t *Server) spawnListener() (net.Listener, error) {
	conf := &tcplisten.Config{
		ReusePort:   true,
		DeferAccept: t.options.DeferAccept,
		FastOpen:    t.options.FastOpen,
	}

	return conf.NewListener("tcp", t.laddr)
}

// canReusePort attempts to test if the OS can use the SO_REUSEPORT socket
// option. It does this by atempting to create a new net.Listener configured
// with that flag set.
func (t *Server) canReusePort() (bool, error) {
	conf := &tcplisten.Config{
		ReusePort: true,
	}

	ln, err := conf.NewListener("tcp", t.laddr)
	if err != nil {
		if strings.Contains(err.Error(), "SO_REUSEPORT") {
			return false, nil
		}

		return false, err // unkown if can use
	}

	ln.Close()
	return true, nil
}

// serveTCP accepts new TCP connections on a given listener and passes them to
// manageTCP. If a error is encountered it is checked to see if it is a
// temporary network error. If it is the routine sleeps for one millisecond
// before continuing, otherwise it sends the error to the provided errorChan
// and returns
func (t *Server) serveTCP(ln net.Listener, errorChan chan error) {
	var errorWait = time.Millisecond

	for {
		// Get next connection and log failures if a logger has been provided
		conn, err := ln.Accept()
		if err != nil {
			// Check if we can recover
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {

				// Temporary error, log, sleep and continue
				if Logger != nil {
					Logger.Printf(
						"A call to listener.Accept failed with error: %v",
						err,
					)
				}

				time.Sleep(errorWait)
				continue
			}

			// Non temporary error
			errorChan <- err
			return
		}

		// Handle the new connection
		go manageTCP(conn, t.handler)
	}
}

func (t *Server) ListenTCP() error {
	reuse, err := t.canReusePort()
	if err != nil {
		if Logger != nil {
			Logger.Printf(
				"Failed to detect SO_REUSEPORT capability error: %v",
				err,
			)
		}
	}

	// We can run multiple listeners (one per CPU) because the OS supports the
	// SO_REUSEPORT socket option
	if reuse {
		if Logger != nil {
			Logger.Printf("Starting paralell listener")
		}

		p := runtime.NumCPU()

		// Create channel with buffer size equal to number of listeners to
		// prevent hanging serveTCP go routines after one has errored
		errorChan := make(chan error, p)

		// Spawn listeners for each CPU and pass them to serveTCP
		for i := 0; i < p; i++ {
			ln, err := t.spawnListener()

			// Log the failure if we have a logger then continue to spawn more
			if err != nil {
				if Logger != nil {
					Logger.Printf(
						"Failed to start listener %v of %v Error: %v",
						i, p, err,
					)
				}

				continue // Try to start some
			}

			// Ensure the listener gets closed
			defer ln.Close()

			// Start the server on this listener
			go t.serveTCP(ln, errorChan)
		}

		// Wait for one of the listeners to errorn then return
		// The defer statements should clear up the other listeners
		return <-errorChan
	}

	// OS does not support the SO_REUSEPORT socket option se we have to use the
	// stdlib single thread listener
	if Logger != nil {
		Logger.Printf("Starting single listener")
	}

	ln, err := net.Listen("tcp", t.laddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	errorChan := make(chan error)
	go t.serveTCP(ln, errorChan)
	return <-errorChan
}
