package sinks

import (
	"github.com/influx6/faux/sink"
)

// Drain emits all entries into nothingness.
type Drain struct{}

// Emit implements the sink.Sink interface and does nothing with the 
// provided entry.
func (Drain) Emit(e sink.Entry) error {
	return nil
}