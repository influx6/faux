package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileFS implements a simple store for storing and retrieving
// file from underneath filesystem.
type FileFS struct {
	Dir string
}

// Persist saves giving file into FileFS.Dir, overwriting any same file
// existing.
func (fs FileFS) Persist(file string, data []byte) error {
	targetPath := filepath.Join(fs.Dir, file)
	targetDir := filepath.Dir(targetPath)

	if _, err := os.Stat(targetDir); err != nil {
		if mkerr := os.MkdirAll(targetDir, 0500); mkerr != nil {
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

// RemoveContents remove FileFS.Dir and contents.
func (fs FileFS) RemoveContents() error {
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

// Retrieve retrieves giving file path within FileFS.Dir.
func (fs FileFS) Retrieve(file string) ([]byte, error) {
	targetFile := filepath.Join(fs.Dir, file)
	if _, err := os.Stat(targetFile); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(targetFile)
}
