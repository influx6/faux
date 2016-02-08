package main

import (
	"github.com/influx6/faux/loop/web"
	"github.com/influx6/faux/vfx"
)

func main() {

	// Initialize the loop engine for vfx.
	vfx.Init(web.Loop)

}
