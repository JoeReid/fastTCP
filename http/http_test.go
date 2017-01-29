package http_test

import (
	"fmt"
	"github.com/JoeReid/fastTCP"
	"github.com/JoeReid/fastTCP/http"
	"io/ioutil"
	goHTTP "net/http"
	"testing"
	"time"
)

func TestHTTP(t *testing.T) {
	hs := http.NewHTTPServer(
		goHTTP.HandlerFunc(func(rw goHTTP.ResponseWriter, r *goHTTP.Request) {
			fmt.Fprint(rw, "Test")
		}),
	)

	server := fastTCP.NewServer(":6543", hs.NewConn, fastTCP.TCPOptions{})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	resp, err := goHTTP.Get("http://localhost:6543/")
	if err != nil {
		t.Log(err)
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}

	t.Logf("%s", buff)
}
