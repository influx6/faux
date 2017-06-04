package http

import (
	"net/http"

	"github.com/influx6/faux/context"
)

const (
	maxsize = 32 << 40

	// MultipartKey defines the key used to store multipart Form.
	MultipartKey = "MultiPartForm"
)

// Params defines a function to return all parameter values and query values
// retrieved from the request.
func Params(r *http.Request, mxsize int64) (context.Context, error) {
	if mxsize <= 0 {
		mxsize = maxsize
	}

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(mxsize); err != nil {
		return nil, err
	}

	ctx := context.New()

	if r.MultipartForm != nil {
		ctx.Set(MultipartKey, r.MultipartForm)
	}

	for name, val := range r.Form {
		ctx.Set(name, val)
	}

	return ctx, nil
}
