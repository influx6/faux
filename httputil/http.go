package httputil

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"path"
	"strings"

	"fmt"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
)

const (
	maxsize = 32 << 40

	// MultipartKey defines the key used to store multipart Form.
	MultipartKey = "MultiPartForm"
)

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
	fn Handler
	m  metrics.Metrics
}

// ServeHTTP implements http.Handler.ServeHttp method.
func (h handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(SetRequest(r), SetResponseWriter(w), SetMetrics(h.m))
	if err := ctx.InitForms(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.fn(ctx)
}

// ServeHandler returns a http.Handler which serves request to the provided Handler.
func ServeHandler(m metrics.Metrics, h Handler) http.Handler {
	return handlerImpl{m: m, fn: h}
}

// GzipServer returns a http.Handler which handles the necessary bits to gzip or ungzip
// file resonses from a http.FileSystem.
func GzipServer(m metrics.Metrics, fs http.FileSystem, gzipped bool) http.Handler {
	return handlerImpl{m: m, fn: GzipServe(fs, gzipped)}
}

// GzipServe returns a Handler which handles the necessary bits to gzip or ungzip
// file resonses from a http.FileSystem.
func GzipServe(fs http.FileSystem, gzipped bool) Handler {
	return func(ctx *Context) error {
		reqURL := path.Clean(ctx.Path())
		if reqURL == "./" || reqURL == "." {
			ctx.Redirect(http.StatusMovedPermanently, "/")
			return nil
		}

		if !strings.HasPrefix(reqURL, "/") {
			reqURL = "/" + reqURL
		}

		file, err := fs.Open(reqURL)
		if err != nil {
			return err
		}

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		if ctx.HasHeader("Accept-Encoding", "gzip") && gzipped {
			ctx.SetHeader("Content-Encoding", "gzip")
			http.ServeContent(ctx.Response(), ctx.Request(), stat.Name(), stat.ModTime(), file)
			return nil
		}

		if ctx.HasHeader("Accept-Encoding", "gzip") && !gzipped {
			ctx.SetHeader("Content-Encoding", "gzip")

			gwriter := gzip.NewWriter(ctx.Response())
			defer gwriter.Close()

			_, err := io.Copy(gwriter, file)
			if err != nil && err != io.EOF {
				return err
			}

			return nil
		}

		if !ctx.HasHeader("Accept-Encoding", "gzip") && gzipped {
			gzreader, err := gzip.NewReader(file)
			if err != nil {
				return err
			}

			var bu bytes.Buffer
			_, err = io.Copy(&bu, gzreader)
			if err != nil && err != io.EOF {
				return err
			}

			http.ServeContent(ctx.Response(), ctx.Request(), stat.Name(), stat.ModTime(), bytes.NewReader(bu.Bytes()))
			return nil
		}

		http.ServeContent(ctx.Response(), ctx.Request(), stat.Name(), stat.ModTime(), file)
		return nil
	}
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
