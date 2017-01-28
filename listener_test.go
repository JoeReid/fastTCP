package fastTCP_test

import (
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

	var recieved = make([]byte, len(testData))
	_, err = conn.Read(recieved)
	if err != nil {
		t.Fatal(err)
	}

	if testData != string(recieved) {
		t.Fatalf("Strings not equal")
	}
}

func TestTCP(t *testing.T) {
	server := fastTCP.NewServer(":6543", testHandler, fastTCP.TCPOptions{})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()

	t.Log("This test is not a load test")
	t.Log("time values may be incorrect")
	t.Log("this routine is too slow to test the load on the network code")

	time.Sleep(time.Second)

	t.Run("TCP one connection", func(t *testing.T) {
		testOnce(t)
	})
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
