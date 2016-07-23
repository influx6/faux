package fs

import (
	"io"
	"os"

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

	WriteBytes([]byte) FileSystem
	WriteWriter(io.Writer) FileSystem

	CreateFile(string) FileSystem
	WriteFile(string) FileSystem

	ReadDir(string) FileSystem
	Mkdir(string) FileSystem

	DeleteFile(string) FileSystem
	DeleteDir(string) FileSystem

	SkipStat(func(os.FileInfo) bool) FileSystem
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
	f.Node = f.Signal(func(ctx pub.Ctx, d interface{}) {

	})
	return f
}

func (f *fs) ReadReader(r io.Reader) FileSystem { return f }

func (f *fs) WriteBytes(data []byte) FileSystem  { return f }
func (f *fs) WriteWriter(w io.Writer) FileSystem { return f }

func (f *fs) CreateFile(path string) FileSystem { return f }
func (f *fs) WriteFile(path string) FileSystem  { return f }

func (f *fs) ReadDir(path string) FileSystem { return f }
func (f *fs) Mkdir(path string) FileSystem   { return f }

func (f *fs) DeleteFile(path string) FileSystem { return f }
func (f *fs) DeleteDir(path string) FileSystem  { return f }

func (f *fs) SkipStat(fn func(os.FileInfo) bool) FileSystem { return f }
func (f *fs) UnwrapStats() FileSystem                       { return f }
