package memory

import (
	"github.com/influx6/faux/sink"
)

// Memory defines a struct which implements a memory collector for sinks.
type Memory struct {
	Data []sink.SentryJSON
}

// Emit adds the giving SentryJSON into the internal slice.
func (m *Memory) Emit(sjn sink.SentryJSON) error {
	m.Data = append(m.Data, sjn)
	return nil
}
