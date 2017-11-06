package httputil

import (
	"bytes"
	"compress/gzip"
	"io"
	"mime"
	"net/http"
	"path"
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

	if err := h.Handler(ctx); err != nil {
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

// GzipServer returns a http.Handler which handles the necessary bits to gzip or ungzip
// file resonses from a http.FileSystem.
func GzipServer(fs http.FileSystem, gzipped bool, mw ...Middleware) http.Handler {
	zipper := GzipServe(fs, gzipped)
	if len(mw) != 0 {
		return handlerImpl{Handler: MWi(mw...)(zipper)}
	}

	return handlerImpl{Handler: zipper}
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

		mime := GetFileMimeType(stat.Name())
		ctx.AddHeader("Content-Type", mime)

		if ctx.HasHeader("Accept-Encoding", "gzip") && gzipped {
			ctx.SetHeader("Content-Encoding", "gzip")
			defer ctx.Status(http.StatusOK)
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

			ctx.Status(http.StatusOK)

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

			defer ctx.Status(http.StatusOK)
			http.ServeContent(ctx.Response(), ctx.Request(), stat.Name(), stat.ModTime(), bytes.NewReader(bu.Bytes()))
			return nil
		}

		defer ctx.Status(http.StatusOK)
		http.ServeContent(ctx.Response(), ctx.Request(), stat.Name(), stat.ModTime(), file)
		return nil
	}
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
