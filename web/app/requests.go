package app

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

//==============================================================================

// ResponseRequest defines a response object which holds the request  object
// associated with it and allows you write out the behaviour.
type ResponseRequest struct {
	// FormParams Param
	Params Param
	ResponseWriter
	R *http.Request
}

// Respond renders out a JSON response and status code giving using the Render
// function.
func (r *ResponseRequest) Respond(code int, data interface{}) {
	Render(code, r.R, r, data)
}

// RespondError renders out a error response into the request object.
func (r *ResponseRequest) RespondError(code int, err error) {
	RenderErrorWithStatus(code, err, r.R, r)
}

//==============================================================================

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
}

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0, 0}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

//==============================================================================
