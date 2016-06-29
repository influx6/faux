package databind

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// StringWriter turns a series of bytes into stringed counterparts using hex notation
type StringWriter struct {
	W io.Writer
}

// Write meets the io.Write interface Write method
func (sw *StringWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	var sod = []byte(`\x00`)
	var b byte

	for n, b = range p {
		sod[2] = lowerhex[b/16]
		sod[3] = lowerhex[b%16]
		sw.W.Write(sod)
	}

	n++

	return
}

// createCompressWriter creates a non-gzip compressed writer
func createUnCompressWriter(w io.Writer) io.WriteCloser {
	return &nopWriter{w}
}

// createCompressWriter creates a gzip compressed writer
func createCompressWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(&StringWriter{W: w})
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reason
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

// ByName Implement sort.Interface for []os.FileInfo based on Name()
type ByName []os.FileInfo

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }

func makeAbsolute(path string) string {
	as, _ := filepath.Abs(path)
	return as
}

// makeRelative removes the if found the beginning slash '/'
func makeRelative(path string) string {
	if path[0] == '/' {
		return path[1:]
	}
	return path
}

func hasIn(paths []string, dt string) bool {
	for _, so := range paths {
		if strings.Contains(so, dt) || so == dt {
			return true
		}
	}
	return false
}

func hasExt(paths []string, dt string) bool {
	for _, so := range paths {
		if so == dt {
			return true
		}
	}
	return false
}

func getDirListings(dir string) ([]os.FileInfo, error) {
	//open up the filepath since its a directory, read and sort
	cdir, err := os.Open(dir)

	if err != nil {
		return nil, err
	}

	defer cdir.Close()

	files, err := cdir.Readdir(0)

	if err != nil {
		return nil, err
	}

	sort.Sort(ByName(files))

	return files, nil
}

type nopWriter struct {
	w io.Writer
}

func (n *nopWriter) Close() error {
	return nil
}

func (n *nopWriter) Write(b []byte) (int, error) {
	return n.w.Write(b)
}
