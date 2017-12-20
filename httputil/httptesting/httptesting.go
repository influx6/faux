package httptesting

import (
	"io"
	"net/http/httptest"

	"net/http"

	"github.com/influx6/faux/httputil"
)

// GET returns a new Context using GET method.
func Get(path string, body io.Reader, res *httptest.ResponseRecorder) *httputil.Context {
	return NewRequest("GET", path, body, res)
}

// Delete returns a new Context using DELETE method.
func Delete(path string, body io.Reader, res *httptest.ResponseRecorder) *httputil.Context {
	return NewRequest("DELETE", path, body, res)
}

// Put returns a new Context using PUT method.
func Put(path string, body io.Reader, res *httptest.ResponseRecorder) *httputil.Context {
	return NewRequest("PUT", path, body, res)
}

// Post returns a new Context using PUT method.
func Post(path string, body io.Reader, res *httptest.ResponseRecorder) *httputil.Context {
	return NewRequest("POST", path, body, res)
}

// Patch returns a new Context using PUT method.
func Patch(path string, body io.Reader, res *httptest.ResponseRecorder) *httputil.Context {
	return NewRequest("PATCH", path, body, res)
}

// NewRequest returns a new instance of a httputil.Context with provided parameters.
func NewRequest(method string, path string, body io.Reader, res http.ResponseWriter) *httputil.Context {
	req := httptest.NewRequest(method, path, body)
	return httputil.NewContext(
		httputil.SetRequest(req),
		httputil.SetResponseWriter(res),
	)
}
