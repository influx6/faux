package falto

import "io"

// lowerhex contains the set of hexadecimal letters and numbers.
const lowerhex = "0123456789abcdef"

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
