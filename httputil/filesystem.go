package httputil

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/influx6/faux/filesystem"
)

// GzipServer returns a http.Handler which handles the necessary bits to gzip or ungzip
// file resonses from a http.FileSystem.
func GzipServer(fs filesystem.FileSystem, gzipped bool, mw ...Middleware) http.Handler {
	zipper := GzipServe(fs, gzipped)
	if len(mw) != 0 {
		return handlerImpl{Handler: MWi(mw...)(zipper)}
	}

	return handlerImpl{Handler: zipper}
}

// GzipServe returns a Handler which handles the necessary bits to gzip or ungzip
// file resonses from a http.FileSystem.
func GzipServe(fs filesystem.FileSystem, gzipped bool) Handler {
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

// HTTPGzipServer returns a http.Handler which handles the necessary bits to gzip or ungzip
// file resonses from a http.FileSystem.
func HTTPGzipServer(fs http.FileSystem, gzipped bool, mw ...Middleware) http.Handler {
	zipper := HTTPGzipServe(fs, gzipped)
	if len(mw) != 0 {
		return handlerImpl{Handler: MWi(mw...)(zipper)}
	}

	return handlerImpl{Handler: zipper}
}

// HTTPGzipServe returns a Handler which handles the necessary bits to gzip or ungzip
// file resonses from a http.FileSystem.
func HTTPGzipServe(fs http.FileSystem, gzipped bool) Handler {
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
