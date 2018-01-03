package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// FilePortal defines an error which exposes methods to
// treat a underline store has a file system.
type FilePortal interface {
	Name() string
	Has(string) bool
	RemoveAll() error
	Remove(string) error
	Dirs() ([]FilePortal, error)
	Save(string, []byte) error
	Get(string) ([]byte, error)
	Within(string) (FilePortal, error)
}

// FileFS implements a simple store for storing and retrieving
// file from underneath filesystem.
type FileFS struct {
	Dir string
}

// Within returns a new FilePortal which exists within the current path.
// It enforces all operations to occur within provided path.
func (fs FileFS) Within(path string) (FilePortal, error) {
	return FileFS{Dir: filepath.Join(fs.Dir, path)}, nil
}

// Name returns underline name of giving FS.
func (fs FileFS) Name() string {
	return filepath.Base(fs.Dir)
}

// Has return true/false if giving file exists in directory of fs.
func (fs FileFS) Has(file string) bool {
	if _, err := os.Stat(filepath.Join(fs.Dir, file)); err != nil {
		return false
	}
	return true
}

// Dirs returns all top level directory in file.
func (fs FileFS) Dirs() ([]FilePortal, error) {
	item, err := os.Open(fs.Dir)
	if err != nil {
		return nil, err
	}

	dirs, err := item.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var portals []FilePortal
	for _, dir := range dirs {
		portals = append(portals, FileFS{
			Dir: filepath.Join(fs.Dir, dir.Name()),
		})
	}

	return portals, nil
}

// Save saves giving file into FileFS.Dir, overwriting any same file
// existing.
func (fs FileFS) Save(file string, data []byte) error {
	targetPath := filepath.Join(fs.Dir, file)
	targetDir := filepath.Dir(targetPath)

	if _, err := os.Stat(targetDir); err != nil {
		if mkerr := os.MkdirAll(targetDir, 0700); mkerr != nil {
			return mkerr
		}
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	defer targetFile.Close()

	_, err = targetFile.Write(data)
	return err
}

// RemoveAll removes all files within FileFS.Dir and contents.
func (fs FileFS) RemoveAll() error {
	if err := os.RemoveAll(fs.Dir); err != nil {
		if perr, ok := err.(*os.PathError); ok && perr.Err == os.ErrNotExist {
			return nil
		}
		return err
	}
	return nil
}

// Remove deletes giving file path within FileFS.Dir.
func (fs FileFS) Remove(file string) error {
	targetFile := filepath.Join(fs.Dir, file)
	if _, err := os.Stat(targetFile); err != nil {
		return nil
	}

	return os.Remove(targetFile)
}

// Get retrieves giving file path within FileFS.Dir.
func (fs FileFS) Get(file string) ([]byte, error) {
	targetFile := filepath.Join(fs.Dir, file)
	if _, err := os.Stat(targetFile); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(targetFile)
}
