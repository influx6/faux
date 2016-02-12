package main

import (
	"time"

	"github.com/influx6/faux/loop/web"
	"github.com/influx6/faux/vfx"
	"github.com/influx6/faux/vfx/animations/boundaries"
)

func main() {

	vfx.Init(web.Loop)

	width := vfx.NewAnimationSequence(".zapps",
		vfx.TimeStat(1*time.Second, "ease-in", true, true, true),
		&boundaries.Width{Width: 500})

	vfx.Animate(width)
}
