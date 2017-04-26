package dom

import "github.com/gopherjs/gopherjs/js"

// Pointer implements the Point interface.
type Pointer struct {
	*js.Object
}

// GetPointAtLength returns a giving point by the provide step value.
func (p *Pointer) GetPointAtLength(step float64) (float64, float64) {
	res := p.Object.Call("getPointAtLength", step)
	if res == js.Undefined || res == nil {
		return 0, 0
	}

	return res.Get("x").Float(), res.Get("y").Float()
}
