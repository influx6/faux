package colors

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"

	"github.com/influx6/faux/utils"
	colorful "github.com/lucasb-eyer/go-colorful"
)

//==============================================================================

// colorReg defines a regexp for matching rgb/rgba header content.
var colorReg = regexp.MustCompile("[rgb|rgba]\\(([\\d\\.,\\s]+)\\)")

// IsRGBFormat returns true/false if the giving string is a rgb/rgba format data.
func IsRGBFormat(c string) bool {
	return colorReg.MatchString(c)
}

// rgbHeader defines a regexp for matching rgb/rgba header content.
var rgbHeader = regexp.MustCompile("rgb\\(([\\d\\.,\\s]+)\\)")
var rgbHeader2 = regexp.MustCompile("rgb\\(([\\d\\.,]+)\\)")

// IsRGB returns true/false if the giving string is a rgb format data.
func IsRGB(c string) bool {
	return rgbHeader.MatchString(c) || rgbHeader2.MatchString(c)
}

// IsHex returns true/false if the giving value is a hex color representation.
func IsHex(c string) bool {
	if strings.HasPrefix(c, "#") {
		if clen := len(c); clen == 3 || clen == 7 {
			return true
		}
	}

	return false
}

// rgbaHeader defines a regexp for matching rgb/rgba header content.
var rgbaHeader = regexp.MustCompile("rgba\\(([\\d\\.,\\s]+)\\)")

// rgbaHeader defines a regexp for matching rgb/rgba header content.
var rgbaHeader2 = regexp.MustCompile("rgba\\(([\\d\\.,]+)\\)")

// IsRGBA returns true/false if the giving string is a rgba format data.
func IsRGBA(c string) bool {
	return rgbaHeader.MatchString(c) || rgbaHeader2.MatchString(c)
}

var hsl = regexp.MustCompile("hsl\\((\\d+),\\s*([\\d.]+)%,\\s*([\\d.]+)%\\)")

// IsHSL returns true/false if the giving string is a rgba format data.
func IsHSL(c string) bool {
	return hsl.MatchString(c)
}

// ParseHSL pulls out the rgb/rgba information from a hsl color format  from the
// provided string.
func ParseHSL(hslData string) (float64, float64, float64) {
	subs := hsl.FindStringSubmatch(hslData)

	h := utils.ParseFloat(subs[1])
	s := utils.ParseFloat(subs[2]) / 100
	l := utils.ParseFloat(subs[3]) / 100

	return h, s, l
}

// HSL2RGB converts color values in hsl to rgb.
func HSL2RGB(h, s, l float64) (int, int, int) {
	if s == 0 {
		return int(l), int(l), int(l)
	}

	var q float64

	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = (l + s) - (l * s)
	}

	p := 2 * (l - q)

	r := Hue(p, q, h+1/3) //* 255
	g := Hue(p, q, h)     //* 255
	b := Hue(p, q, h-1/3) //* 255

	return int(r), int(g), int(b)
}

// PastelMix takes the provided colorbases and returns a random ix of
// random colors.
func PastelMix(r, g, b int64) (int64, int64, int64) {
	nr := rand.Int63n(256)
	gr := rand.Int63n(256)
	br := rand.Int63n(256)

	nr = (nr + r) / 2
	gr = (gr + g) / 2
	br = (gr + g) / 2

	return nr, gr, br
}

// DistinctMix returns a distinct color value based on using rainbow non-overlapping
// values which can easily be used for distinct items where color is important.
func DistinctMix() (int, int, int) {
	hue := math.Floor(rand.Float64()*30) * 12
	return HSL2RGB(hue, 0.9, 0.6)
}

// Hue takes the provided values and returns a hue value.
func Hue(p, q, t float64) float64 {
	if t < 0 {
		t++
	}

	if t > 1 {
		t--
	}

	if t < 1/6 {
		return p + (q-p)*6*t
	}

	if t < 1/2 {
		return q
	}

	if t < 2/3 {
		return p + (q-p)*(2/3-t)*6
	}

	return p
}

// ParseRGB pulls out the rgb/rgba information from a rgba(9,9,9,9) type
// formatted string.
func ParseRGB(rgbData string) (int, int, int, float64) {
	subs := colorReg.FindStringSubmatch(rgbData)

	if len(subs) < 2 {
		return 0, 0, 0, 0
	}

	rc := strings.Split(subs[1], ",")

	var r, g, b int
	var alpha float64

	r = utils.ParseInt(rc[0])
	g = utils.ParseInt(rc[1])
	b = utils.ParseInt(rc[2])

	if len(rc) > 3 {
		alpha = utils.ParseFloat(rc[3])
	} else {
		alpha = 1
	}

	return r, g, b, alpha
}

// HexToRGB turns a hexademicmal color into rgba format.
// Returns the read, green and blue values as int.
func HexToRGB(hex string) (int, int, int) {
	c, err := colorful.Hex(hex)
	if err != nil {
		return 0, 0, 0
	}

	r, g, b, _ := c.RGBA()
	return int(r), int(g), int(b)
}

// RGBToHSL converts the provided rgb value to HSL format.
func RGBToHSL(r, g, b int) (float64, float64, float64) {
	c := colorful.Color{
		R: float64(r) / 255,
		B: float64(g) / 255,
		G: float64(b) / 255,
	}

	return c.Hsl()
}

// HexToRGBA turns a hexademicmal color into rgba format.
// Alpha values ranges from 0-100
func HexToRGBA(hex string, alpha int) string {
	r, g, b := HexToRGB(hex)
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r, g, b, float64(alpha)/100)
}

// doubleString doubles the giving string.
func doubleString(c string) string {
	return fmt.Sprintf("%s%s", c, c)
}
