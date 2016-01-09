package styles

import (
	"strconv"

	"github.com/influx6/faux/domtrees"
)

// Size presents a basic stringifed unit
type Size string

// Px returns the value in "%px" format
func Px(pixels int) Size {
	return Size(strconv.Itoa(pixels) + "px")
}

// Color provides the color style value
func Color(value string) *domtrees.Style {
	return &domtrees.Style{Name: "color", Value: value}
}

// Height provides the height style value
func Height(size Size) *domtrees.Style {
	return &domtrees.Style{Name: "height", Value: string(size)}
}

// FontSize provides the margin style value
func FontSize(size Size) *domtrees.Style {
	return &domtrees.Style{Name: "font-size", Value: string(size)}
}

// Padding provides the margin style value
func Padding(size Size) *domtrees.Style {
	return &domtrees.Style{Name: "padding", Value: string(size)}
}

// Margin provides the margin style value
func Margin(size Size) *domtrees.Style {
	return &domtrees.Style{Name: "margin", Value: string(size)}
}

// Width provides the width style value
func Width(size Size) *domtrees.Style {
	return &domtrees.Style{Name: "width", Value: string(size)}
}
