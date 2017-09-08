package file

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/influx6/faux/metrics"
)

// File defines a struct which implements a memory collector for metricss.
type File struct {
	wl   sync.Mutex
	file *os.File
	path string
}

// New returns a new instance of a File sentry.
func New(path string) *File {
	var fm File
	fm.path = path
	return &fm
}

// Emit adds the giving SentryJSON into the internal slice.
func (f *File) Emit(sjn metrics.SentryJSON) error {
	f.wl.Lock()
	defer f.wl.Unlock()

	if f.file == nil {
		fm, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		f.file = fm
	}

	if err := json.NewEncoder(f.file).Encode(&sjn); err != nil {
		f.file.Close()
		f.file = nil
		return err
	}

	if _, err := f.file.Write([]byte("\r\n")); err != nil {
		f.file.Close()
		f.file = nil
		return err
	}

	if err := f.file.Sync(); err != nil {
		f.file.Close()
		f.file = nil
		return err
	}

	return nil
}
