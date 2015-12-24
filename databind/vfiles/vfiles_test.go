package vfiles

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/influx6/flux"
)

func TestCompressedVirtualFile(t *testing.T) {
	vf := NewVFile("./", "assets/bucklock.txt", "buklock.txt", 30, true, true, func(v *VFile) ([]byte, error) {
		return readData(v, []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x52\x0e\xcb\xcc\xe5\x02\x04\x00\x00\xff\xff\xec\xfe\xa5\xd9\x05\x00\x00\x00"))
	})

	if vf.Size() < 30 {
		flux.FatalFailed(t, "Incorrect size value,expected %d got %d", 30, vf.Size())
	}

	if data, err := vf.Data(); err != nil {
		flux.FatalFailed(t, "Error occured retrieving content: %s", err)
	} else {
		if string(data) != "#Vim\n" {
			flux.FatalFailed(t, "Error in file content expected %s got %s", "shop", data)
		}
	}

	flux.LogPassed(t, "Successfully read contents of virtual file")
}

func TestDebugVirtualFile(t *testing.T) {
	vf := NewVFile("./", "assets/vim.md", "vim.md", 5, true, true, func(v *VFile) ([]byte, error) {
		return readFile(v)
	})

	if vf.Size() > 5 {
		flux.FatalFailed(t, "Incorrect size value,expected %d got %d", 5, vf.Size())
	}

	if data, err := vf.Data(); err != nil {
		flux.FatalFailed(t, "Error occured retrieving content: %s", err)
	} else {
		if string(data) != "#Vim\n" {
			flux.FatalFailed(t, "Error in file content expected %s got %s", "#Vim", data)
		}
	}

	flux.LogPassed(t, "Successfully read contents of virtual file")
}

func TestPlainVirtualFile(t *testing.T) {
	vf := NewVFile("./", "assets/bucklock.txt", "buklock.txt", 5, true, true, func(v *VFile) ([]byte, error) {
		return []byte("shop"), nil
	})

	if vf.Size() > 5 {
		flux.FatalFailed(t, "Incorrect size value,expected %d got %d", 5, vf.Size())
	}

	if data, err := vf.Data(); err != nil {
		flux.FatalFailed(t, "Error occured retrieving content: %s", err)
	} else {
		if string(data) != "shop" {
			flux.FatalFailed(t, "Error in file content expected %s got %s", "shop", data)
		}
	}

	flux.LogPassed(t, "Successfully read contents of virtual file")
}

func TestVirtualDir(t *testing.T) {
	var root = NewDirCollector()

	root.Set("assets", func() *VDir {
		var dir = NewVDir("assets", ".", "./", true)

		dir.AddDirectory("tests", func() *VDir {
			return root.Get("assets/tests")
		})

		dir.AddFile(NewVFile("./", "assets/shop.md", "shop.md", 5, true, false, func(v *VFile) ([]byte, error) {
			return []byte("shop"), nil
		}))
		return dir
	}())

	root.Set("assets/tests", func() *VDir {
		var dir = NewVDir("assets/tests", "./tests", "./", false)
		dir.AddFile(NewVFile("./", "assets/tests/lock.md", "tests/lock.md", 5, true, false, func(v *VFile) ([]byte, error) {
			return []byte("shop"), nil
		}))
		return dir
	}())

	if _, err := root.GetDir("/assets/"); err != nil {
		flux.FatalFailed(t, "Unable to located asset dir in dirCollection")
	}

	if _, err := root.GetDir("/assets/tests/"); err != nil {
		flux.FatalFailed(t, "Unable to located asset/tests dir in dirCollection")
	}

	dir, _ := root.GetDir("/assets/")

	_, err := dir.GetFile("shop.md")

	if err != nil {
		flux.FatalFailed(t, "Unable to located assets/shop.md file: %s", err)
	}

	to, err := dir.GetDir("/tests/")

	if err != nil {
		flux.FatalFailed(t, "Unable to located asset/tests directory: %s", err)
	}

	_, err = to.GetFile("lock.md")

	if err != nil {
		flux.FatalFailed(t, "Unable to located asset/tests/lock.md directory: %s", err)
	}

	al, err := root.GetFile("/assets/tests/lock.md")

	if data, err := al.Data(); err != nil {
		flux.FatalFailed(t, "failed to load assets/tests/lock.md contents: %s", err)
	} else if string(data) != ("shop") {
		flux.FatalFailed(t, "incorrect assets/tests/lock.md content, expected %q got %q", "shop", data)
	}

	if err != nil {
		flux.FatalFailed(t, "Unable to located asset/tests/lock.md directory: %s", err)
	}
}

func readFile(v *VFile) ([]byte, error) {
	fo, err := ioutil.ReadFile(v.RealPath())
	if err != nil {
		return nil, fmt.Errorf("---> assets.readFile: Error reading file: %s at %s: %s\n", v.Name(), v.RealPath(), err)
	}
	return fo, nil
}

func readData(v *VFile, data []byte) ([]byte, error) {
	if !v.Decompress {
		return readVData(v, data)
	}

	return readEData(v, data)
}
