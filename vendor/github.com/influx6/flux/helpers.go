package flux

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//ToDuration returns the duration in millisecond of a given time from the current time
func ToDuration(ms time.Time) time.Duration {
	return Elapsed(ms, time.Now())
}

//Elapsed takes a time and delta it from the current time to return a duration in milliseconds
func Elapsed(ms time.Time, diff time.Time) time.Duration {
	return Elapsed(ms, diff) * time.Millisecond
}

//ElapsedIn returns the elapsed time in int64
func ElapsedIn(ms time.Time, diff time.Time) int64 {
	return diff.UTC().Sub(ms.UTC()).Nanoseconds() / 1e6
}

//FileCloser provides a means of closing a file
type FileCloser struct {
	*os.File
	path string
}

//Close ends and deletes the file
func (f *FileCloser) Close() error {
	ec := f.File.Close()
	log.Printf("Will Remove %s", f.path)
	ex := os.Remove(f.path)

	if ex == nil {
		return ec
	}

	return ex
}

//NewFileCloser returns a new file closer
func NewFileCloser(path string) (*FileCloser, error) {
	ff, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return &FileCloser{ff, path}, nil
}

//BufferCloser closes a byte.Buffer
type BufferCloser struct {
	*bytes.Buffer
}

//NewBufferCloser returns a new closer for a bytes.Buffer
func NewBufferCloser(bu *bytes.Buffer) *BufferCloser {
	return &BufferCloser{bu}
}

//Close resets the internal buffer
func (b *BufferCloser) Close() error {
	b.Buffer.Reset()
	return nil
}

//GzipWalker walks a path and turns it into a tar written into a bytes.Buffer
func GzipWalker(file string, tmp io.Writer) error {
	f, err := os.Open(file)

	if err != nil {
		return err
	}

	defer f.Close()

	//gzipper
	gz := gzip.NewWriter(tmp)
	defer gz.Close()

	_, err = io.Copy(gz, f)

	return err
}

//TarWalker walks a path and turns it into a tar written into a bytes.Buffer
func TarWalker(rootpath string, w io.Writer) error {
	tz := tar.NewWriter(w)
	defer tz.Close()

	walkFn := func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}

		np, err := filepath.Rel(rootpath, path)
		if err != nil {
			return err
		}

		fl, err := os.Open(path)
		if err != nil {
			return err
		}

		defer fl.Close()

		var h *tar.Header
		if h, err = tar.FileInfoHeader(info, ""); err != nil {
			return err
		}

		h.Name = np

		if err := tz.WriteHeader(h); err != nil {
			return err
		}

		if _, err := io.Copy(tz, fl); err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(rootpath, walkFn)
	if err != nil {
		return err
	}

	return nil
}

//Backwards takes a value and walks Backward till 0
func Backwards(to int, fx func(int)) {
	for i := to; i > 0; i-- {
		fx(i)
	}
}

//Forwards takes a value and walks Backward till 0
func Forwards(to int, fx func(int)) {
	for i := 1; i <= to; i++ {
		fx(i)
	}
}

//BackwardsIf takes a value and walks Backward till 0 unless the stop function is called
func BackwardsIf(to int, fx func(int, func())) {
	state := true
	for i := to; i > 0; i-- {
		if !state {
			break
		}
		fx(i, func() { state = false })
	}
}

//ForwardsIf takes a value and walks Backward till 0 unless the stop func is called
func ForwardsIf(to int, fx func(int, func())) {
	state := true
	for i := 1; i <= to; i++ {
		if !state {
			break
		}
		fx(i, func() { state = false })
	}
}

//BackwardsSkip takes a value and walks Backward till 0 unless the skip function is called it will go through all sequence
func BackwardsSkip(to int, fx func(int, func())) {
	for i := to; i > 0; i-- {
		fx(i, func() { i-- })
	}
}

//ForwardsSkip takes a value and walks Backward till 0 unless the skip func is called it will go throuh all sequence
func ForwardsSkip(to int, fx func(int, func())) {
	for i := 1; i <= to; i++ {
		fx(i, func() { i++ })
	}
}

//Report provides a nice abstaction for doing basic report
func Report(e error, msg string) {
	if e != nil {
		log.Fatalf("Message: (%s) with Error: (%s)", msg, e.Error())
	} else {
		log.Printf("Message: (%s) with NoError", msg)
	}
}

// SilentRecoveryHandler provides a recovery handler functions for use to automate the recovery processes
func SilentRecoveryHandler(tag string, opFunc func() error) error {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	if err := opFunc(); err != nil {
		return err
	}

	return nil
}

//RecoveryHandler provides a recovery handler functions for use to automate the recovery processes
func RecoveryHandler(tag string, opFunc func() error) error {
	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 10000)
			count := runtime.Stack(trace, true)
			trace = trace[:count]
			log.Printf("---------%s-Panic----------------:", strings.ToUpper(tag))
			log.Printf("Error: %+s", err)
			log.Printf("Stack of %d bytes: %+s\n", count, trace)
			log.Printf("---------%s--END-----------------:", strings.ToUpper(tag))
		}
	}()

	if err := opFunc(); err != nil {
		return err
	}

	return nil
}

// SilentRecoveryHandlerCallback provides a recovery handler functions for use to automate the recovery processes
func SilentRecoveryHandlerCallback(tag string, opFunc func() error, cb func(interface{})) error {
	defer func() {
		if err := recover(); err != nil {
			cb(err)
			return
		}
	}()

	if err := opFunc(); err != nil {
		return err
	}

	return nil
}

//RecoveryHandlerCallback provides a recovery handler functions for use to automate the recovery processes
func RecoveryHandlerCallback(tag string, opFunc func() error, cb func(interface{})) error {
	defer func() {
		if err := recover(); err != nil {
			cb(err)
			trace := make([]byte, 10000)
			count := runtime.Stack(trace, true)
			log.Printf("---------%s-Panic----------------:", strings.ToUpper(tag))
			log.Printf("Error: %+s", err)
			log.Printf("Stack of %d bytes: %+s\n", count, trace)
			log.Printf("---------%s--END-----------------:", strings.ToUpper(tag))
		}
	}()

	if err := opFunc(); err != nil {
		return err
	}

	return nil
}

//GoDefer letsw you run a function inside a goroutine that gets a defer recovery
func GoDefer(title string, fx func()) {
	go RecoveryHandler(title, func() error {
		fx()
		return nil
	})
}

//GoSilent letsw you run a function inside a goroutine that gets a defer recovery
func GoSilent(title string, fx func()) {
	go SilentRecoveryHandler(title, func() error {
		fx()
		return nil
	})
}

//GoDeferCall letsw you run a function inside a goroutine that gets a defer recovery
func GoDeferCall(title string, fx func(), cb func(interface{})) {
	go RecoveryHandlerCallback(title, func() error {
		fx()
		return nil
	}, cb)
}

//GoSilentCall lets you run a function inside a goroutine that gets a defer recovery
func GoSilentCall(title string, fx func(), cb func(interface{})) {
	go SilentRecoveryHandlerCallback(title, func() error {
		fx()
		return nil
	}, cb)
}

//Close provides a basic io.WriteCloser write method
func (w *FuncWriter) Close() error {
	w.fx = nil
	return nil
}

//Write provides a basic io.Writer write method
func (w *FuncWriter) Write(b []byte) (int, error) {
	w.fx(b)
	return len(b), nil
}

//NewFuncWriter returns a new function writer instance
func NewFuncWriter(fx func([]byte)) *FuncWriter {
	return &FuncWriter{fx}
}

type (
	//FuncWriter provides a means of creation io.Writer on functions
	FuncWriter struct {
		fx func([]byte)
	}
)
