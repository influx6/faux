package fs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/influx6/faux/pub"
)

// FileSystem provides a core interface for FileSystem operations, where each
// instance of this type is series of operation for a single task or objective.
// This allows you to construct a nice chain of events, that can be recalled
// over and over again to do the same thing purely and nothing else.
type FileSystem interface {
	pub.Node

	ReadFile(string) FileSystem
	ReadReader(io.Reader) FileSystem
	ReplayBytes([]byte) FileSystem
	ReplayReader(io.Reader) FileSystem

	WriteBytes([]byte) FileSystem
	WriteWriterBytes([]byte) FileSystem
	WriteWriter(io.Writer) FileSystem

	CreateFile(string) FileSystem
	CloseFile(string) FileSystem
	WriteFile(string) FileSystem

	Mkdir(string, bool) FileSystem
	ReadDir(string) FileSystem
	WalkDir(string) FileSystem

	Remove(string) FileSystem
	RemoveAll(string) FileSystem

	SkipStat(func(ExtendedFileInfo) bool) FileSystem
	UnwrapStats() FileSystem
}

// New returns a new FileSystem
func New() FileSystem {
	var fu fs
	fu.Node = pub.InverseMagic(pub.IdentityHandler())
	return &fu
}

//==============================================================================

type fs struct {
	pub.Node
}

// newFS returns a new instance pointer of a fs struct, which attaches the
// inverse of the provided Node to itself. The inverse ensures this is a
// serially chained link.
func newFS(node pub.Node) *fs {
	var fu fs
	fu.Node = node.Inverse()
	return &fu
}

// ReadFile adds a readFile operation whoes contents get passed to the next
// event/Node/Task in the link.
func (f *fs) ReadFile(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, _ interface{}) {
		file, err := os.Open(path)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		defer file.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, file)
		if err != nil && err != io.EOF {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, buf.Bytes())
	})
	return f
}

// ReadReader reads the data pulled from the reader everytime it gets called.
func (f *fs) ReadReader(r io.Reader) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, _ interface{}) {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil && err != io.EOF {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, buf.Bytes())
	})
	return f
}

// ReplayBytes resends the data provided everytime it is called.
func (f *fs) ReplayBytes(b []byte) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, _ interface{}) {
		ctx.RW().Write(ctx, b)
	})
	return f
}

// ReplayReader reads the data pulled from the reader everytime, buffers it
// and returns data everytime it gets called.
func (f *fs) ReplayReader(r io.Reader) FileSystem {
	var buf bytes.Buffer
	var read bool

	f.Node = f.MustSignal(func(ctx pub.Ctx, _ interface{}) {
		if read {
			ctx.RW().Write(ctx, buf.Bytes())
			return
		}

		_, err := io.Copy(&buf, r)
		if err != nil && err != io.EOF {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, buf.Bytes())
	})
	return f
}

// WriteBytes writes the giving bytes to a path it expects to receive when called,
// It appends the provided data to that path continously.
// It passes the data passed in to the its subscribers to
// both allow the chain of events to continue and to allow others to use the data
// as they please.
func (f *fs) WriteBytes(data []byte) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, path string) {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0500)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		defer file.Close()

		written, err := file.Write(data)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		if written != len(data) {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, data)
	})
	return f
}

// WriteWriterBytes writes the giving bytes to a giving writer it receives when
// the previous node sends out its data. It appends the provided data to that
// path continously.
// It passes the data passed in to the its subscribers to
// both allow the chain of events to continue and to allow others to use the data
// as they please.
func (f *fs) WriteWriterBytes(data []byte) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, w io.Writer) {
		written, err := w.Write(data)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		if written != len(data) {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, w)
	})
	return f
}

// WriteWriter expects to recieve []bytes as input and writes the provided
// bytes into the writer it recieves as argument. It returns error if the total
// written does not match the size of the bytes. It passes the incoming data
// down the pipeline.
func (f *fs) WriteWriter(w io.Writer) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, data []byte) {
		written, err := w.Write(data)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		if written != len(data) {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, data)
	})
	return f
}

// CreateFile creates the giving file within the provided path and returns the
// file object down its pipeline. The file it creates is always in create mode,
// hence over-written any file by the same name at the provided path.
// If it gets a []byte slice as its argument. It will write the provided slice
// into the file.
func (f *fs) CreateFile(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, data []byte) {
		file, err := os.OpenFile(path, os.O_CREATE, 0500)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		if len(data) > 0 {
			_, err := file.Write(data)
			if err != nil {
				ctx.RW().Write(ctx, err)
				return
			}
		}

		ctx.RW().Write(ctx, file)
	})
	return f
}

// CloseFile expects to receive a file object which it calls the close function
// on. It passes true if the file closed without error.
func (f *fs) CloseFile(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, f *os.File) {
		if err := f.Close(); err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, true)
	})
	return f
}

// WriteFile either creates or opens an existing file for appending. It passes
// the file object for this files down its pipeline.
// If it gets a []byte slice as its argument. It will write the provided slice
// into the file.
func (f *fs) WriteFile(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, data []byte) {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0500)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		if len(data) > 0 {
			_, err := file.Write(data)
			if err != nil {
				ctx.RW().Write(ctx, err)
				return
			}
		}

		ctx.RW().Write(ctx, file)
	})

	return f
}

// ExtendedFileInfo composes a os.FileInfo to provide the fullPath property
// for a giving fileInfo.
type ExtendedFileInfo interface {
	os.FileInfo
	Path() string
	Dir() string
}

// NewExtendFileInfo returns a structure which implements the ExtendedFileInfo
// interface.
func NewExtendFileInfo(info os.FileInfo, root string) ExtendedFileInfo {
	ef := extendFileInfo{
		FileInfo: info,
		path:     filepath.Join(root, info.Name()),
		root:     root,
	}

	return ef
}

type extendFileInfo struct {
	os.FileInfo
	path string
	root string
}

// Dir returns the directory of the provided file.
func (e extendFileInfo) Dir() string {
	return e.root
}

// Path returns the path of the provided file.
func (e extendFileInfo) Path() string {
	return e.path
}

// ReadDir reads the giving path if indeed is a directory, else passing down
// an error down the provided pipeline. It extends the provided os.FileInfo
// with a structure that implements the ExtendFileInfo interface. It sends the
// individual fileInfo instead of the slice of FileInfos.
func (f *fs) ReadDir(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, data []byte) {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0500)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		defer file.Close()

		dirs, err := file.Readdir(-1)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		for _, dir := range dirs {
			ctx.RW().Write(ctx, NewExtendFileInfo(dir, path))
		}

		ctx.RW().WriteEnd(ctx)
	})
	return f
}

// WalkDir walks the giving path if indeed is a directory, else passing down
// an error down the provided pipeline. It extends the provided os.FileInfo
// with a structure that implements the ExtendFileInfo interface. It sends the
// individual fileInfo instead of the slice of FileInfos.
func (f *fs) WalkDir(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, _ interface{}) {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0500)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		defer file.Close()

		dirs, err := file.Readdir(-1)
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		for _, dir := range dirs {
			dirInfo := NewExtendFileInfo(dir, path)
			ctx.RW().Write(ctx, dirInfo)

			// If this is a sysmbol link, then continue we won't read through it.
			if _, err := os.Readlink(dirInfo.Dir()); err == nil {
				continue
			}

			ctx.RW().Read(dirInfo, ctx.Ctx())
		}

		ctx.RW().WriteEnd(ctx)
	})
	return f
}

// Mkdir creates a directly returning the path down the pipeline. If the chain
// flag is on, then mkdir when it's pipeline receives a non-empty string as
// an argument, will join the string recieved with the path provided.
// This allows chaining mkdir paths down the pipeline.
func (f *fs) Mkdir(path string, chain bool) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, root string) {
		if chain && path != "" {
			path = filepath.Join(root, path)
		}

		if err := os.MkdirAll(path, 0500); err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, path)
	})
	return f
}

// Remove deletes the giving path and passes the path down
// the pipeline.
func (f *fs) Remove(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, data []byte) {
		if err := os.Remove(path); err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, path)
	})
	return f
}

// RemoveAll deletes the giving path and its subpaths if its a directory
// and passes the path down
// the pipeline.
func (f *fs) RemoveAll(path string) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, data []byte) {
		if err := os.Remove(path); err != nil {
			ctx.RW().Write(ctx, err)
			return
		}

		ctx.RW().Write(ctx, path)
	})
	return f
}

// SkipStat takes a function to filter out the FileInfo that are running through
// its pipeline. This allows you to define specific file paths you wish to treat.
// If the filter function returns true, then any FileInfo/ExtendFileInfo that
// match its criteria are sent down its pipeline.
func (f *fs) SkipStat(filter func(ExtendedFileInfo) bool) FileSystem {
	f.Node = f.MustSignal(func(ctx pub.Ctx, info ExtendedFileInfo) {
		if filter(info) {
			ctx.RW().Write(ctx, info)
		}
	}, false, true)
	return f
}

// UnwrapStats takes the provided ExtendFileInfo and unwraps them into their
// full path, allows you to retrieve the strings path.
func (f *fs) UnwrapStats() FileSystem {
	var pack []string
	f.Node = f.SignalEnd(func(ctx pub.Ctx) {
		ctx.RW().Write(ctx, pack)
		pack = nil
	}).MustSignal(func(ctx pub.Ctx, info ExtendedFileInfo) {
		pack = append(pack, info.Path())
	}, false, true)
	return f
}
