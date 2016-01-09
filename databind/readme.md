#DataBind
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/influx6/faux/binddata)

Provides a convenient set of tools for handling template files and turning assets into embeddable go files

##Example

  - Embedding

    ```go

    	bf, err := NewBindFS(&BindFSConfig{
    		InDir:   "./",
    		OutDir:     "./tests/debug",
    		Package: "debug",
    		File:    "debug",
    		Gzipped: true,
            NoDecompression: true,
            Production: false,
    	})

    	if err != nil {
             panic("directory path is not valid")
    	}

      //to get this to create and embed the files,simple call .Record()
    	err = bf.Record() // you can call this as many times as you want to update a go file

      // A generated file called `debug.go` will exists in ./tests/debug/


    ```

    - Loading a generated asset file

    ```go

      import (
        "github.com/influx6/assets/tests/debug"
        "net/http"
      )

      func main(){

        //to retrieve a directory,simply do:
        fixtures,err := debug.RootDirectory.GetDir("/fixtures/")

        //to retrieve a file,simply do:
        basic,err := debug.RootDirectory.GetFile("/fixtures/base/basic.tmpl")

        // create a http.FileServer from the global RootDirectory listing
        rootFs := http.FileServer(debug.RootDirectory)

        // or use the root VirtualDirectory as a http.FileSystem
        rootFs2 := http.FileServer(debug.RootDirectory.Root())

        //or use any sub-directory you want
        fixturesFs := http.FileServer(debug.RootDirectory.Get("/fixtures/"))

      }
    ```
