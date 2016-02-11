// +build ignore

package main

import (
	"path/filepath"

	"github.com/influx6/faux/databind"
)

var root = "../."
var tests = filepath.Join(root, "tests")

func main() {
	debugBindFS()
	productionBindFS()

	debugCompressBindFS()
	productionNoDecompressBindFS()
}

func debugCompressBindFS() {
	bf, err := databind.NewBindFS(&databind.BindFSConfig{
		InDir:           root,
		OutDir:          filepath.Join(tests, "debugnodecompress"),
		Package:         "debug",
		File:            "debug",
		Gzipped:         true,
		NoDecompression: true,
	})

	err = bf.Record()
	if err != nil {
		panic(err)
	}
}

func debugBindFS() {
	bf, err := databind.NewBindFS(&databind.BindFSConfig{
		InDir:   root,
		OutDir:  filepath.Join(tests, "debug"),
		Package: "debug",
		File:    "debug",
		Gzipped: false,
	})

	err = bf.Record()
	if err != nil {
		panic(err)
	}
}

func productionBindFS() {
	bf, err := databind.NewBindFS(&databind.BindFSConfig{
		InDir:      root,
		OutDir:     filepath.Join(tests, "prod"),
		Package:    "prod",
		File:       "prod",
		Gzipped:    true,
		Production: true,
	})

	err = bf.Record()
	if err != nil {
		panic(err)
	}
}

func productionNoDecompressBindFS() {
	bf, err := databind.NewBindFS(&databind.BindFSConfig{
		InDir:           root,
		OutDir:          filepath.Join(tests, "prodnodecompress"),
		Package:         "prod",
		File:            "prod",
		Gzipped:         true,
		NoDecompression: true,
		Production:      true,
	})

	err = bf.Record()
	if err != nil {
		panic(err)
	}
}
