package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
)

type HTTPServer struct {
	handler http.Handler
}

func NewHTTPServer(h http.Handler) *HTTPServer {
	return &HTTPServer{h}
}

func (h *HTTPServer) NewConn(tcp io.ReadWriter) {
	br := bufio.NewReader(tcp)

	req, err := http.ReadRequest(br)
	if err != nil {
		return
	}

	rw := newFastResponseWriter(req, tcp)
	h.handler.ServeHTTP(rw, req)
}

func newFastResponseWriter(req *http.Request, conn io.Writer) http.ResponseWriter {
	return &fastResponseWriter{
		response: &http.Response{
			Header:           http.Header(make(map[string][]string)),
			Proto:            req.Proto,
			ProtoMajor:       req.ProtoMajor,
			ProtoMinor:       req.ProtoMinor,
			TransferEncoding: req.TransferEncoding,
		},
		conn: conn,
	}
}

type fastBody struct {
	buff   *bytes.Buffer
	closed bool
}

func (f *fastBody) Read(b []byte) (int, error) {
	if f.closed {
		return 0, errors.New("Closed")
	}
	return f.buff.Read(b)
}

func (f *fastBody) Close() error {
	if f.closed {
		return errors.New("Error already closed")
	}

	f.closed = true
	return nil
}

type fastResponseWriter struct {
	response *http.Response
	conn     io.Writer
}

func (f *fastResponseWriter) Header() http.Header {
	return f.response.Header
}

func (f *fastResponseWriter) Write(bdy []byte) (int, error) {
	if f.response.StatusCode == 0 {
		f.WriteHeader(http.StatusOK)
	}

	f.response.ContentLength = int64(len(bdy))
	f.response.Body = &fastBody{buff: bytes.NewBuffer(bdy)}

	b, err := httputil.DumpResponse(f.response, true)
	if err != nil {
		return 0, err
	}

	_, err = f.conn.Write(b)
	if err != nil {
		return 0, err
	}

	return len(bdy), nil
}

func (f *fastResponseWriter) WriteHeader(h int) {
	f.response.StatusCode = h
	f.response.Status = http.StatusText(h)
}
