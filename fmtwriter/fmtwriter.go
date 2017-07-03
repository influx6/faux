package fmtwriter

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/influx6/faux/process"
)

// WriterTo defines a takes the contents of a provided io.WriterTo
// against go fmt and returns the result.
type WriterTo struct {
	io.WriterTo
	goimport        bool
	attemptFallback bool
}

// New returns a new instance of FmtWriterTo.
func New(wt io.WriterTo, useGoImports bool, attemptFallbackToFmtInError bool) *WriterTo {
	return &WriterTo{WriterTo: wt, goimport: useGoImports, attemptFallback: attemptFallbackToFmtInError}
}

// WriteTo writes the content of the source after running against gofmt to the
// provider writer.
func (fm WriterTo) WriteTo(w io.Writer) (int64, error) {
	var backinput, input, inout, inerr bytes.Buffer

	if n, err := fm.WriterTo.WriteTo(io.MultiWriter(&input, &backinput)); err != nil && err != io.EOF {
		return n, err
	}

	cmdName := "gofmt"

	if fm.goimport {
		cmdName = "goimports"
	}

	cmd := process.Command{
		Name:  cmdName,
		Level: process.RedAlert,
	}

	if err := cmd.Run(context.Background(), &inout, &inerr, &input); err != nil {

		// If we must attempt to fallback to gofmt, due to goimport error, attempt to
		if fm.goimport && fm.attemptFallback {
			return (WriterTo{WriterTo: &backinput}).WriteTo(w)
		}

		errcount, _ := inerr.WriteTo(w)
		linecount, _ := fmt.Fprintf(w, "\n-----------------------\n")
		outcount, _ := backinput.WriteTo(w)

		return (errcount + int64(linecount) + outcount), fmt.Errorf("GoFmt Error: %+q (See generated file for fmt Error)", err)
	}

	return inout.WriteTo(w)
}
