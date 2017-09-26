package memory

import (
	"github.com/influx6/faux/metrics"
)

// Memory defines a struct which implements a memory collector for metricss.
type Memory struct {
	Data []metrics.Entry
}

// Emit adds the giving SentryJSON into the internal slice.
func (m *Memory) Emit(en metrics.Entry) error {
	m.Data = append(m.Data, en)
	return nil
}
