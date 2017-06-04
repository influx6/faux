package httputil

import (
	"net/http"
	"strings"

	"github.com/influx6/faux/context"
)

const (
	maxsize = 32 << 40

	// MultipartKey defines the key used to store multipart Form.
	MultipartKey = "MultiPartForm"
)

// Params defines a function to return all parameter values and query values
// retrieved from the request.
func Params(ctx context.Context, r *http.Request, multipartFormSize int64) error {
	if multipartFormSize <= 0 {
		multipartFormSize = maxsize
	}

	switch {
	case strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data"):
		if err := r.ParseMultipartForm(multipartFormSize); err != nil {
			return err
		}

		break
	default:
		if err := r.ParseForm(); err != nil {
			return err
		}
	}

	if r.MultipartForm != nil {
		ctx.Set(MultipartKey, r.MultipartForm)
	}

	for name, val := range r.Form {
		ctx.Set(name, val)
	}

	return nil
}
