package main

import (
	"fmt"

	"github.com/influx6/faux/vfx"
)

func main() {

	div := vfx.Document().QuerySelector("div.offset")

	top, left := vfx.Offset(div)
	ptop, pleft := vfx.Position(div)

	div.SetInnerHTML(fmt.Sprintf("Offset: %.2f %.2f\nPosition: %.2f %.2f\n", top, left, ptop, pleft))
}
