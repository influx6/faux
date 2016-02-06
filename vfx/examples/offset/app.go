package main

import (
	"fmt"

	"github.com/influx6/faux/vfx"
)

func main() {

	div := vfx.Document().QuerySelector("div.offset")

	top, left := vfx.Offset(div)
	ptop, pleft := vfx.Position(div)

	var color = "#ff32CC"

	div.SetInnerHTML(fmt.Sprintf(`
    Offset: %.2f %.2f
    Position: %.2f %.2f
    Color: Hex(%s) Rgba(%s)
  `, top, left, ptop, pleft, color, vfx.RGBA(color, 50)))
}
