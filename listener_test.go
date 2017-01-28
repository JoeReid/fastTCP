package fastTCP_test

import (
	"fmt"
	"github.com/JoeReid/fastTCP"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func testHandler(f io.ReadWriter) {
	io.Copy(f, f)
}

func testOnce(t *testing.T) {
	conn, err := net.Dial("tcp", ":6543")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	testData := "FOOBAR_foobar foobar"
	_, err = conn.Write([]byte(testData))
	if err != nil {
		t.Fatal(err)
	}

	var received = make([]byte, len(testData))
	_, err = conn.Read(received)
	if err != nil {
		t.Fatal(err)
	}

	if testData != string(received) {
		t.Fatalf("Strings not equal")
	}
}

type serverLogger struct {
	t *testing.T
}

func (s *serverLogger) Printf(format string, v ...interface{}) {
	s.t.Log(fmt.Sprintf(format, v...))
}

func (s *serverLogger) Println(v ...interface{}) {
	s.t.Log(v...)
}

func TestServerStop(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	for i := 0; i < 5; i++ {
		server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{})
		go func() {
			err := server.ListenTCP()
			if err != nil {
				panic(err)
			}
		}()

		time.Sleep(50 * time.Millisecond)

		server.Stop()
		time.Sleep(time.Millisecond)
	}
}

func TestTCPDefault(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func TestTCPDefaultIPv6(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{
		IPv6: true,
	})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func TestTCPDeferAccept(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{
		DeferAccept: true,
	})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func TestTCPDeferAcceptIPv6(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{
		DeferAccept: true,
		IPv6:        true,
	})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func TestTCPFastOpen(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{
		FastOpen: true,
	})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func TestTCPFastOpenIPv6(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{
		FastOpen: true,
		IPv6:     true,
	})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func TestTCPDeferAcceptFastOpen(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{
		DeferAccept: true,
		FastOpen:    true,
	})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func TestTCPDeferAcceptFastOpenIPv6(t *testing.T) {
	fastTCP.Logger = &serverLogger{t}

	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{
		DeferAccept: true,
		FastOpen:    true,
		IPv6:        true,
	})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(50 * time.Millisecond)

	t.Run("TCP one connection", testMultiple(1))
	t.Run("TCP ten connection", testMultiple(10))
	t.Run("TCP hundred connection", testMultiple(100))
}

func testMultiple(mul int) func(*testing.T) {
	return func(t *testing.T) {
		wg := &sync.WaitGroup{}
		for i := 0; i < mul; i++ {
			wg.Add(1)
			go func() {
				testOnce(t)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
