package fs

import "github.com/influx6/faux/pub"

// FileSystem provides a core interface for FS operations.
type FileSystem interface {
	pub.Node
	ReadFile(string) FS
	ReadDir(string) FS
	Write() FS
}

// FS defines a global handle which exposes the fs API.
var FS fs

type fs struct {
	pub.Node
}

// ReadFile returns a new FS which recieves contents from the giving file.
func (f fs) ReadFile(path string) FS {

}
