package main

import (
	"fmt"
	"time"

	"github.com/influx6/faux/loop/web"
	"github.com/influx6/faux/vfx"
	"github.com/influx6/faux/vfx/animations/boundaries"
)

func main() {

	vfx.Init(web.Loop)

	width := vfx.NewAnimationSequence(".zapps",
		vfx.TimeStat(500*time.Millisecond, "ease-in", true, false, true),
		&boundaries.Width{Width: 500})

	fmt.Printf("width sequence: %+s\n", width)
}
