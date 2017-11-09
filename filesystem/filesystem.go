package filesystem

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// errors ...
var (
	ErrExist       = errors.New("Path exists")
	ErrNotExist    = errors.New("Path does not exists")
	ErrToManyParts = errors.New("Prefix must contain only two parts: /bob/log")
)

// FileMode defines the mode value of a file.
type FileMode uint32

// FileInfo defines an interface representing a file into.
type FileInfo interface {
	Name() string
	Size() int64
	ModTime() time.Time
	IsDir() bool
}

// File defines a interface for representing a file.
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]FileInfo, error)
	Stat() (FileInfo, error)
}

// FileSystem defines a interface for a virutal filesystem
// representing a file.
type FileSystem interface {
	Open(string) (File, error)
}

// New returns a new instance of a VirtualFileSystem as a FileSystem type.
func New(fn GetFile) FileSystem {
	return VirtualFileSystem{
		GetFileFunc: fn,
	}
}

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

// IsDir returns false because this is a virtual file.
func (vf *VirtualFile) IsDir() bool {
	return false
}

// Mode returns associated file mode of file.
func (vf *VirtualFile) Mode() FileMode {
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
func (vf *VirtualFile) Stat() (FileInfo, error) {
	return vf, nil
}

// Readdir returns nil slices as this is a file not a directory.
func (vf *VirtualFile) Readdir(n int) ([]FileInfo, error) {
	return nil, nil
}

// Close returns nothing.
func (vf *VirtualFile) Close() error {
	return nil
}

// GetFile define a function type that returns a VirtualFile type
// or an error.
type GetFile func(string) (File, error)

// StripPrefix returns a new GetFile which wraps the previous provided
// GetFile and always strips provided prefix from incoming path.
func StripPrefix(prefix string, from GetFile) GetFile {
	return func(path string) (File, error) {
		return from(strings.TrimPrefix(path, prefix))
	}
}

// VirtualFileSystem connects a series of functions which are provided
// to retrieve bundled files and serve to a http server. It implements
// http.FileSystem interface.
type VirtualFileSystem struct {
	GetFileFunc GetFile
}

// Open returns associated file with given name if found else
// returning an error. It implements http.FileSystem.Open method.
func (v VirtualFileSystem) Open(name string) (File, error) {
	if v.GetFileFunc == nil {
		return nil, ErrNotExist
	}

	vfile, err := v.GetFileFunc(name)
	if err != nil {
		return nil, err
	}

	return vfile, nil
}

type systemNode struct {
	prefix string
	root   FileSystem
	nodes  map[string]FileSystem
}

// SystemGroup allows the combination of multiple filesystem to
// respond to incoming request based on initial path prefix.
type SystemGroup struct {
	ml      sync.Mutex
	systems map[string]systemNode
}

// NewSystemGroup returns a new instance of SystemGroup.
func NewSystemGroup() *SystemGroup {
	return &SystemGroup{
		systems: make(map[string]systemNode),
	}
}

// MustRegister will panic if the prefix and FileSystem failed to register
// It returns itself if successfully to allow chaining.
func (fs *SystemGroup) MustRegister(prefix string, m FileSystem) *SystemGroup {
	if err := fs.Register(prefix, m); err != nil {
		panic(err)
	}

	return fs
}

// Register adds giving file system to handling paths with given prefix.
func (fs *SystemGroup) Register(prefix string, m FileSystem) error {
	defer fs.ml.Unlock()
	fs.ml.Lock()

	if strings.Contains(prefix, "\\") {
		prefix = filepath.ToSlash(prefix)
	}

	prefix = strings.TrimPrefix(prefix, "/")
	prefix = strings.TrimSuffix(prefix, "/")

	parts := strings.Split(prefix, "/")
	if len(parts) > 2 {
		return ErrToManyParts
	}

	root := parts[0]
	if len(parts) == 1 {
		if node, ok := fs.systems[root]; ok {
			node.root = m
			fs.systems[root] = node
			return nil
		}

		var node systemNode
		node.prefix = root
		node.root = m
		node.nodes = make(map[string]FileSystem)
		fs.systems[root] = node
		return nil
	}

	tail := parts[1]
	if node, ok := fs.systems[root]; ok {
		node.nodes[tail] = m
		fs.systems[root] = node
		return nil
	}

	var node systemNode
	node.prefix = root
	node.nodes = make(map[string]FileSystem)
	node.nodes[tail] = m
	fs.systems[root] = node

	return nil
}

// Open attempts to locate giving path from different file systems, else returning error.
func (fs *SystemGroup) Open(path string) (File, error) {
	if len(path) == 0 {
		return nil, ErrNotExist
	}

	defer fs.ml.Unlock()
	fs.ml.Lock()

	if strings.Contains(path, "\\") {
		path = filepath.ToSlash(path)
	}

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")
	if node, ok := fs.systems[parts[0]]; ok {
		if len(parts) > 1 {
			if subnode, ok := node.nodes[parts[1]]; ok {
				return subnode.Open(filepath.Join(parts[2:]...))
			}
		}

		if node.root != nil {
			return node.root.Open(filepath.Join(parts[1:]...))
		}
	}

	return nil, ErrNotExist
}
