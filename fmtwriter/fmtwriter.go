package fmtwriter

import (
	"bytes"
	"context"
	"io"

	"github.com/influx6/faux/process"
)

// WriterTo defines a takes the contents of a provided io.WriterTo
// against go fmt and returns the result.
type WriterTo struct {
	io.WriterTo
}

// New returns a new instance of FmtWriterTo.
func New(wt io.WriterTo) *WriterTo {
	return &WriterTo{WriterTo: wt}
}

// WriteTo writes the content of the source after running against gofmt to the
// provider writer.
func (fm WriterTo) WriteTo(w io.Writer) (int64, error) {
	cmd := process.Command{
		Name:  "gofmt",
		Level: process.RedAlert,
	}

	var input bytes.Buffer

	if n, err := fm.WriterTo.WriteTo(&input); err != nil && err != io.EOF {
		return n, err
	}

	if err := cmd.Run(context.Background(), w, w, &input); err != nil {
		return 0, nil
	}

	return int64(input.Len()), nil
}
