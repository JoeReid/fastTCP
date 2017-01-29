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

func handleOK(w goHTTP.ResponseWriter, r *goHTTP.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "Status OK")
}

func handleNotFound(w goHTTP.ResponseWriter, r *goHTTP.Request) {
	w.WriteHeader(404)
	fmt.Fprint(w, "Status Not Found")
}

func handleHeader(k, v string) func(w goHTTP.ResponseWriter, r *goHTTP.Request) {
	return func(w goHTTP.ResponseWriter, r *goHTTP.Request) {
		w.WriteHeader(200)
		w.Header().Set(k, v)
		fmt.Fprint(w, "Status OK, with custom headers")
	}
}

func TestHTTPOK(t *testing.T) {
	hs := http.NewHTTPServer(goHTTP.HandlerFunc(handleOK))

	server := fastTCP.NewServer(":6544", hs.NewConn, fastTCP.TCPOptions{})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	time.Sleep(50 * time.Millisecond)

	resp, err := goHTTP.Get("http://localhost:6544/")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200OK got %v", resp.StatusCode)
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", buff)
}

func TestHTTPNotFound(t *testing.T) {
	hs := http.NewHTTPServer(goHTTP.HandlerFunc(handleNotFound))

	server := fastTCP.NewServer(":6544", hs.NewConn, fastTCP.TCPOptions{})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	time.Sleep(50 * time.Millisecond)

	resp, err := goHTTP.Get("http://localhost:6544/")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Expected 404 got %v", resp.StatusCode)
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", buff)
}

func TestHTTPHeader(t *testing.T) {
	hs := http.NewHTTPServer(goHTTP.HandlerFunc(handleHeader("TestHeader", "FooBar")))

	server := fastTCP.NewServer(":6544", hs.NewConn, fastTCP.TCPOptions{})
	go func() {
		err := server.ListenTCP()
		if err != nil {
			panic(err)
		}
	}()
	defer server.Stop()

	time.Sleep(50 * time.Millisecond)

	resp, err := goHTTP.Get("http://localhost:6544/")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200OK got %v", resp.StatusCode)
	}

	if resp.Header.Get("TestHeader") != "FooBar" {
		t.Fatalf(
			"Expected header 'TestHeader' to be set to 'FooBar' but got '%v'",
			resp.Header.Get("TestHeader"),
		)
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", buff)
}
