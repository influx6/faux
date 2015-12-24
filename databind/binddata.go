package databind

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
)

var vfile = filepath.Join(os.Getenv("GOPATH"), "src/github.com/influx6/faux/databind/vfiles/vfiles.go")

// DevelopmentMode represents development mode for bfs files
const DevelopmentMode = 0

// ProductionMode repesents a production assembly mode for bfs
const ProductionMode = 1

// BindFSConfig provides a configuration struct for BindFS
type BindFSConfig struct {
	InDir           string        //directory path use as source
	OutDir          string        //directory path to save file in
	Package         string        //package name for the file
	File            string        //file name of the file
	Gzipped         bool          // to enable gzipping of filecontents
	NoDecompression bool          // active only when Gzipped is true,this disables decompression of data response or forces compression of output when in debug mode
	Production      bool          // to enable production mode as default
	ValidPath       PathValidator //use to filter allowed paths
	Mux             PathMux       //use to mutate path look
	Ignore          *regexp.Regexp
}

// BindFS provides the struct for creating and updating a go file containing static assets from a directory
type BindFS struct {
	config       *BindFSConfig
	listing      *DirListing
	mode         int64
	endpoint     string
	endpointDir  string
	inputDir     string
	curDir       string
	vfileContent string
}

// NewBindFS returns a new BindFS instance or an error if it fails to located directory
func NewBindFS(config *BindFSConfig) (*BindFS, error) {
	vali := config.ValidPath
	mux := config.Mux

	pwd, _ := os.Getwd()
	input := filepath.Join(pwd, config.InDir)
	endpoint := filepath.Join(pwd, config.OutDir, config.File+".go")
	endpointDir := filepath.Dir(endpoint)

	(config).ValidPath = func(path string, in os.FileInfo) bool {
		if strings.Contains(path, ".git") {
			return false
		}

		if config.Ignore != nil && config.Ignore.MatchString(path) {
			return false
		}

		if path == config.OutDir || strings.Contains(path, config.OutDir) {
			return false
		}

		// check if this contains the endpoint directory. If so, ignore.
		if strings.Contains(path, endpointDir) || strings.Contains(filepath.Join(pwd, path), endpointDir) {
			return false
		}

		if vali != nil {
			return vali(path, in)
		}

		return true
	}

	//clean out pathway or use custom muxer if provided
	(config).Mux = func(path string, in os.FileInfo) string {
		if mux != nil {
			return mux(path, in)
		}
		if filepath.Clean(path) == "." {
			return "/"
		}

		// fmt.Printf("path: %s clean: %s \n", path, filepath.Clean(path))

		return path
	}

	ls, err := DirListings(config.InDir, config.ValidPath, config.Mux)

	if err != nil {
		return nil, fmt.Errorf("---> BindFS: Unable to create dirlisting for %s -> %s", config.InDir, err)
	}

	bf := BindFS{
		config:      config,
		listing:     ls,
		endpoint:    endpoint,
		endpointDir: endpointDir,
		inputDir:    input,
	}

	if config.Production {
		bf.ProductionMode()
	}

	return &bf, nil
}

// Mode returns the current mode of the BindFS
func (bfs *BindFS) Mode() int {
	return int(atomic.LoadInt64(&bfs.mode))
}

// DevMode switches BindFS operations into dev mode
func (bfs *BindFS) DevMode() {
	atomic.StoreInt64(&bfs.mode, DevelopmentMode)
}

// ProductionMode switches BindFS operations into production mode and embeds file data into their corresponding vfiles
func (bfs *BindFS) ProductionMode() {
	atomic.StoreInt64(&bfs.mode, ProductionMode)
}

// Record dumps all the files and dir listings with their corresponding data into a go file within the specified path
func (bfs *BindFS) Record() error {
	bfs.listing.Reload()

	pwd, _ := os.Getwd()
	input := bfs.inputDir
	endpoint := bfs.endpoint
	endpointDir := bfs.endpointDir
	pkgHeader := fmt.Sprintf(packageDetails, bfs.config.Package, input, bfs.config.Package)

	err := os.MkdirAll(endpointDir, 0700)

	if err != nil && err != os.ErrExist {
		return err
	}

	//load the vfile content if not catched.
	if bfs.vfileContent == "" {
		vf, err := os.Open(vfile)
		if err != nil {
			return err
		}

		var vbuf bytes.Buffer
		_, err = io.Copy(&vbuf, vf)
		if err != nil && err != io.EOF {
			vf.Close()
			return err
		}

		// Close the file.
		vf.Close()

		bfs.vfileContent = strings.Replace(vbuf.String(), "package vfiles", "", -1)

	}

	//remove the file for safety and to reduce bloated ouput if file was added in list
	// os.Remove(endpoint)

	boutput, err := os.Create(endpoint)

	if err != nil && err != os.ErrExist {
		return err
	}

	defer boutput.Close()

	var output = bufio.NewWriter(boutput)

	//writes the library package header
	fmt.Fprint(output, pkgHeader)

	//writing the libraries core
	fmt.Fprint(output, bfs.vfileContent)
	fmt.Fprint(output, rootDir)

	var noCompressed bool
	if bfs.Mode() > 0 {
		if bfs.config.Gzipped && bfs.config.NoDecompression {
			noCompressed = true
		} else {
			noCompressed = false
		}
	}

	fmt.Fprint(output, fmt.Sprintf(comFunc, noCompressed))

	// log.Printf("tree: %s", bfs.listing.Listings.Tree)

	//go through the directories listings
	bfs.listing.EachDir(func(dir *BasicAssetTree, path string) {
		// log.Printf("walking dir: %s", path)

		path = filepath.ToSlash(filepath.Clean(path))
		modDir := filepath.ToSlash(filepath.Clean(dir.ModDir))
		pathDir := filepath.ToSlash(filepath.Clean(dir.Dir))
		pathAbs := filepath.ToSlash(filepath.Clean(dir.AbsDir))

		if path == ".." {
			path = "/"
		}

		// if it has a .. at the beginning, remove it.
		if strings.HasPrefix(path, "..") {
			path = strings.TrimPrefix(path, "..")
		}

		if modDir == ".." {
			modDir = "/"
		}

		// if it has a .. at the beginning, remove it.
		if strings.HasPrefix(modDir, "..") {
			modDir = strings.TrimPrefix(modDir, "..")
		}

		//fill up the directory content
		dirContent := fmt.Sprintf(dirRegister, path, modDir, pathDir, pathAbs, dir.root)

		var subs []string
		var data []string

		// go through the subdirectories list and added them
		dir.EachChild(func(child *BasicAssetTree) {
			childDir := filepath.ToSlash(filepath.Clean(child.ModDir))
			baseChildDir := filepath.Base(childDir)

			if baseChildDir == ".." {
				baseChildDir = "/"
			}

			// if it has a .. at the beginning, remove it.
			// if strings.HasPrefix(baseChildDir, "..") {
			// 	baseChildDir = strings.TrimPrefix(baseChildDir, "..")
			// }

			if strings.HasPrefix(childDir, "..") {
				childDir = strings.TrimPrefix(childDir, "..")
			}

			// if it contains the output skip
			//add the sub-directories
			subs = append(subs, fmt.Sprintf(subRegister, baseChildDir, childDir))
		})

		//loadup the files
		dir.Tree.Each(func(modded, real string) {
			modded = filepath.ToSlash(filepath.Clean(modded))
			real = filepath.ToSlash(filepath.Clean(real))
			cleanPwd := filepath.ToSlash(filepath.Clean(pwd))

			// if it has a .. at the beginning, remove it.
			if strings.HasPrefix(modded, "..") {
				modded = strings.TrimPrefix(modded, "..")
			}

			var output string
			if bfs.Mode() == DevelopmentMode {
				stat, _ := os.Stat(filepath.Join(pwd, real))
				var filreadFunc = fileRead
				var size int64

				if stat != nil {
					size = stat.Size()
				}

				if bfs.config.Gzipped && bfs.config.NoDecompression {
					filreadFunc = comfileRead
				}

				output = fmt.Sprintf(debugFile, cleanPwd, modded, real, size, bfs.config.Gzipped, !bfs.config.NoDecompression, filreadFunc)
			} else {
				//production mode is active,we need to load the file contents

				file, err := os.Open(real)

				if err != nil {
					fmt.Printf("---> BindFS.error: failed to loadup %s file -> %s", real, err)
					return
				}

				var data bytes.Buffer
				var writer io.WriteCloser

				if bfs.config.Gzipped {
					writer = createCompressWriter(&data)
				} else {
					writer = createUnCompressWriter(&data)
				}

				n, _ := io.Copy(writer, file)
				file.Close()
				writer.Close()

				var bu []byte

				if bfs.config.Gzipped {
					bu = sanitize(data.Bytes())
				} else {
					bu = data.Bytes()
				}

				var format string
				if bfs.config.Gzipped {
					stringed := fmt.Sprintf("%q", bu)
					stringed = strings.Replace(stringed, `\\`, `\`, -1)
					format = fmt.Sprintf(prodRead, stringed)
				} else {
					format = fmt.Sprintf(prodRead, fmt.Sprintf("`%s`", bu))
				}

				output = fmt.Sprintf(debugFile, cleanPwd, modded, real, n, bfs.config.Gzipped, !bfs.config.NoDecompression, format)
			}

			data = append(data, output)
		})

		dirContent = strings.Replace(dirContent, "{{ subs }}", strings.Join(subs, "\n"), -1)
		dirContent = strings.Replace(dirContent, "{{ files }}", strings.Join(data, "\n"), -1)

		fmt.Fprint(output, fmt.Sprintf(rootInit, dirContent))
	})

	// io.Copy(boutput, output)
	// log.Printf("flushing to file")
	output.Flush()

	return nil
}
