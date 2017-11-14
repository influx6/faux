package httputil

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/influx6/faux/context"
)

const (
	maxsize = 32 << 40

	// MultipartKey defines the key used to store multipart Form.
	MultipartKey = "MultiPartForm"
)

// HTTPError defines custom error that can be used to specify
// status code and message.
type HTTPError struct {
	Code int
	Err  error
}

// Error returns error string. Implements error interface.
func (h HTTPError) Error() string {
	return h.Err.Error()
}

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// handlerImpl implements http.Handler interface.
type handlerImpl struct {
	Handler
}

// ServeHTTP implements http.Handler.ServeHttp method.
func (h handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(SetRequest(r), SetResponseWriter(w))
	if err := ctx.InitForms(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer ctx.ClearFlashMessages()

	if err := h.Handler(ctx); err != nil {
		if httperr, ok := err.(HTTPError); ok {
			http.Error(w, httperr.Error(), httperr.Code)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// HTTPFunc returns a http.HandleFunc which wraps the Handler for usage
// with a server.
func HTTPFunc(nx Handler, befores ...func()) http.HandlerFunc {
	return handlerImpl{Handler: func(ctx *Context) error {
		ctx.response.beforeFuncs = append(ctx.response.beforeFuncs, befores...)
		return nx(ctx)
	}}.ServeHTTP
}

// ServeHandler returns a http.Handler which serves request to the provided Handler.
func ServeHandler(h Handler) http.Handler {
	return handlerImpl{Handler: h}
}

// GetFileMimeType returns associated mime type for giving file extension.
func GetFileMimeType(path string) string {
	ext := filepath.Ext(path)
	extVal := mime.TypeByExtension(ext)
	if extVal == "" {
		extVal = mediaTypes[ext]
	}
	return extVal
}

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

// ParseAuthorization returns the scheme and token of the Authorization string
// if it's valid.
func ParseAuthorization(val string) (authType string, token string, err error) {
	authSplit := strings.SplitN(val, " ", 2)
	if len(authSplit) != 2 {
		err = errors.New("Invalid Authorization: Expected content: `AuthType Token`")
		return
	}

	authType = strings.TrimSpace(authSplit[0])
	token = strings.TrimSpace(authSplit[1])

	return
}
