package boundaries

import (
	"fmt"

	"github.com/influx6/faux/vfx"
)

//==============================================================================

// WidthCSSWriter defines a DeferWriters for writing width properties.
type WidthCSSWriter struct {
	width    int
	unit     string
	priority bool
	elem     vfx.Elemental
}

// Write writes out the necessary output for a css with property.
func (w *WidthCSSWriter) Write() {
	val := fmt.Sprintf("%d%s", w.width, w.unit)
	w.elem.Write("width", val, w.priority)
	w.elem.Sync()
}

//==============================================================================

// Width provides animation sequencing for width properties.
type Width struct {
	Width int
}

// Init returns the initial writers for the sequence.
func (w *Width) Init(stats vfx.Stats, elems vfx.Elementals) vfx.DeferWriters {
	var writers vfx.DeferWriters

	for _, elem := range elems {
		width, priority, _ := elem.Read("width")
		writers = append(writers, &WidthCSSWriter{
			width:    vfx.ParseInt(width),
			unit:     "px",
			priority: priority,
			elem:     elem,
		})
	}

	return writers
}

// Next returns the writers for the current sequence iteration.
func (w *Width) Next(stats vfx.Stats, elems vfx.Elementals) vfx.DeferWriters {
	var writers vfx.DeferWriters

	for _, elem := range elems {
		width, priority, _ := elem.Read("width")

		realWidth := vfx.ParseInt(width)
		change := w.Width - realWidth

		newWidth := easeIn(stats.CurrentIteration(), realWidth, change, stats.TotalIterations())

		fmt.Printf("Width %d-> %d --> %d\n", realWidth, change, newWidth)

		writers = append(writers, &WidthCSSWriter{
			width:    newWidth,
			unit:     "px",
			priority: priority,
			elem:     elem,
		})
	}

	return writers
}

//==============================================================================

// easeInQuad returns a easing value forthe current sequence.
func easeIn(startTime, currentValue, changeInValue, totalTime int) int {
	ms := float64(startTime) / float64(totalTime)
	fmt.Printf("CurrentTime: %d/%d -> %.2f\n", startTime, totalTime, ms)
	cm := float64(changeInValue) * ms
	fmt.Printf("change: %.2f\n", cm)
	return int(cm*ms) + currentValue
}
