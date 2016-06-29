# Pattern
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/influx6/faux/pattern)

  Pattern provides a simple URI pattern matcher, useful for constructing url
  matchers.

## Example

  ```go

  package main

  import "github.com/influx6/faux/pattern"
  import "fmt"

  func main(){

  	r := pattern.New(`/name/{id:[\d+]}/`)

  	params, state := r.Validate(`/name/12/d`)
    if !state {
      panic("No match found")
    }

    fmt.Printf("URL Params: %+s",params)
  }

  ```
