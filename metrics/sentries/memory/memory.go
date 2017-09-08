package memory

import (
	"github.com/influx6/faux/metrics"
)

// Memory defines a struct which implements a memory collector for metricss.
type Memory struct {
	Data []metrics.SentryJSON
}

// Emit adds the giving SentryJSON into the internal slice.
func (m *Memory) Emit(sjn metrics.SentryJSON) error {
	m.Data = append(m.Data, sjn)
	return nil
}
