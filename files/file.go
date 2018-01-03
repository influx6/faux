package files

import (
	"bytes"
	"os"
	"time"
)

// DodFile defines a standin structure for a giving http.File object.
type DodFile struct {
	Reader *bytes.Reader
	Info   os.FileInfo
	Dirs   []os.FileInfo
}

// New returns a new DodFile instance.
func New(name string, size int64, body *bytes.Reader, dirs ...os.FileInfo) *DodFile {
	return &DodFile{
		Dirs:   dirs,
		Reader: body,
		Info: &DodInfo{
			FileName: name,
			FileSize: size,
		},
	}
}

// Read reads data into the byte slice.
func (d *DodFile) Read(b []byte) (int, error) {
	return d.Reader.Read(b)
}

// Seek implements the io.Seeker interface.
func (d *DodFile) Seek(offset int64, whence int) (int64, error) {
	return d.Reader.Seek(offset, whence)
}

// Close returns nil always.
func (d DodFile) Close() error {
	return nil
}

// Stat returns the FileInfo associated wih the Dodfile.
func (d DodFile) Stat() (os.FileInfo, error) {
	return d.Info, nil
}

// Readdir returns the Dirs field as the list of directories.
func (d DodFile) Readdir(dm int) ([]os.FileInfo, error) {
	return d.Dirs, nil
}

// DodInfo implements the FileInfo for usage as a os.FileInfo.
type DodInfo struct {
	FileName string
	FileSize int64
}

// Sys returns always nil.
func (m DodInfo) Sys() interface{} {
	return nil
}

// Mode returns the current time.
func (m DodInfo) Mode() os.FileMode {
	return 0777
}

// ModTime returns the current time.
func (m DodInfo) ModTime() time.Time {
	return time.Now()
}

// Size returns the size associated for the structure.
func (m DodInfo) Size() int64 {
	return m.FileSize
}

// IsDir returns always false.
func (m DodInfo) IsDir() bool {
	return false
}

// IsRegular reports whether m describes a regular file.
func (m DodInfo) IsRegular() bool {
	return true
}

// Perm returns the Unix permission bits in m.
func (m DodInfo) Perm() os.FileMode {
	return 0700
}

// Name returns the name associated for the DodInfo
func (m DodInfo) Name() string {
	return m.FileName
}
