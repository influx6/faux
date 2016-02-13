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
		vfx.TimeStat(vfx.StatConfig{
			Duration: 1 * time.Second,
			Delay:    2 * time.Second,
			Easing:   "ease-in",
			Loop:     4,
			Reverse:  true,
			Optimize: true,
		}),
		&boundaries.Width{Width: 500})

	vfx.Animate(width)
}
