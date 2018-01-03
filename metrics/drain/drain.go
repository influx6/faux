package drain

import (
	"github.com/influx6/faux/metrics"
)

// Drain emits all entries into nothingness.
type Drain struct{}

// Handle implements the metrics.metrics interface and does nothing with the
// provided entry.
func (Drain) Handle(e metrics.Entry) error {
	return nil
}
