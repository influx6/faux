package httputil

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/influx6/faux/context"
)

const (
	maxsize = 32 << 40

	// MultipartKey defines the key used to store multipart Form.
	MultipartKey = "MultiPartForm"
)

// Params defines a function to return all parameter values and query values
// retrieved from the request.
func Params(ctx context.ValueBag, r *http.Request, multipartFormSize int64) error {
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
		if len(val) == 1 {
			ctx.Set(name, val[0])
			continue
		}

		ctx.Set(name, val)
	}

	return nil
}

// ErrorMessage returns a string which contains a json value of a
// given error message to be delivered.
func ErrorMessage(status int, header string, err error) string {
	return fmt.Sprintf(`{
		"status": %d,
		"title": %+q,
		"message": %+q,
	}`, status, header, err)
}

// WriteErrorMessage writes the giving error message to the provided writer.
func WriteErrorMessage(w http.ResponseWriter, status int, header string, err error) {
	http.Error(w, ErrorMessage(status, header, err), status)
}
