package httputil

import (
	"bytes"
	"net/http"
	"os"
	"time"
)

// VirtualFile exposes a slice of []byte and associated name as
// a http.File. It implements http.File interface.
type VirtualFile struct {
	*bytes.Reader
	FileName string
	FileSize int64
	FileMod  time.Time
}

// NewVirtualFile returns a new instance of VirtualFile.
func NewVirtualFile(r *bytes.Reader, filename string, size int64, mod time.Time) *VirtualFile {
	return &VirtualFile{
		Reader:   r,
		FileSize: size,
		FileName: filename,
		FileMod:  mod,
	}
}

// ModTime returns associated mode time for file.
func (vf *VirtualFile) ModTime() time.Time {
	if vf.FileMod.IsZero() {
		return time.Now()
	}

	return vf.FileMod
}

// Sys returns nil has underlying data source.
func (vf *VirtualFile) Sys() interface{} {
	return nil
}

// IsDir returns false because this is a virtual file.
func (vf *VirtualFile) IsDir() bool {
	return false
}

// Mode returns associated file mode of file.
func (vf *VirtualFile) Mode() os.FileMode {
	return 0700
}

// Size returns data file size.
func (vf *VirtualFile) Size() int64 {
	return vf.FileSize
}

// Name returns filename of giving file, either as a absolute or
// relative path.
func (vf *VirtualFile) Name() string {
	return vf.FileName
}

// Stat returns VirtualFile which implements os.FileInfo for virtual
// file to meet http.File interface.
func (vf *VirtualFile) Stat() (os.FileInfo, error) {
	return vf, nil
}

// Readdir returns nil slices as this is a file not a directory.
func (vf *VirtualFile) Readdir(n int) ([]os.FileInfo, error) {
	return nil, nil
}

// Close returns nothing.
func (vf *VirtualFile) Close() error {
	return nil
}

// GetFile define a function type that returns a VirtualFile type
// or an error.
type GetFile func(string) (VirtualFile, error)

// VirtualFileSystem connects a series of functions which are provided
// to retrieve bundled files and serve to a http server. It implements
// http.FileSystem interface.
type VirtualFileSystem struct {
	GetFileFunc GetFile
}

// Open returns associated file with given name if found else
// returning an error. It implements http.FileSystem.Open method.
func (v VirtualFileSystem) Open(name string) (http.File, error) {
	if v.GetFileFunc == nil {
		return nil, os.ErrNotExist
	}

	vfile, err := v.GetFileFunc(name)
	if err != nil {
		return nil, err
	}

	return &vfile, nil
}
