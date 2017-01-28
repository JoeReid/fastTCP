package fastTCP

import (
	"io"
	"net"
	"runtime"
	"sync/atomic"
)

// activeConnections is the number of tcp connections currently open. This
// value is atomicly incremented and decremented when new tcp connections are
// handled by the manageTCP function
var activeConnections int32 = 0

// osThreads is the number of OS threads the applcation should be able to use,
// see runtime.GOMAXPROCS() for clarification on this behaviour.
var osThreads int = runtime.NumCPU()

// This needs to be done for go version < 1.5 but is redundant otherwise
func init() {
	runtime.GOMAXPROCS(osThreads)
}

// updateThreads is used to ensure there are enough osThreads for the network
// conns to block without causing starvation.
func updateThreads(conCount int) {
	if conCount >= osThreads {
		// grow at log n changes in thread count for n connections
		osThreads = osThreads * 2

		runtime.GOMAXPROCS(osThreads)
	}
}

// manageTCP takes a net.Conn and exposes the inner os.File to the provided
// handler function to allow for blocking IO, avoiding go's netpoller.
// To avoid starvation of other go routines during blocking IO, this method
// calls runtime.LockOSThread() and manages the number of available OS threads
// via a growing runtime.GOMAXPROCS() value.
//
// Changes to the number of threads happen at log n rate where n is the number
// of connections by using a double on capacity reached method similar to slice
// growth.
func manageTCP(conn net.Conn, handler func(io.ReadWriter)) {
	// Increment the activeConnections counter and defer decrement
	conCount := atomic.AddInt32(&activeConnections, 1)
	defer atomic.AddInt32(&activeConnections, -1)

	// Check there are enough OS threads then lock the thread in preperation for
	//blocking IO
	updateThreads(int(conCount))
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Get the underlying os file, this call sets the net conn blocking
	file, err := conn.(*net.TCPConn).File()
	conn.Close()

	if err != nil {
		if Logger != nil {
			Logger.Printf(
				"A call to (*net.TCPConn).File() failed with error: %v",
				err,
			)
		}
		return
	}

	defer file.Close()
	handler(file)
}
