package flux

import (
	"fmt"
	"log"
	"testing"
	"time"
)

type caller func(interface{}, int)

type sim struct {
	fx caller
}

func (s *sim) Work(c interface{}, id int) {
	if s.fx != nil {
		s.fx(c, id)
	}
}

func TestWorkPool(t *testing.T) {

	wo, err := NewPool("test.com", "test", &PoolConfig{
		MaxWorkers: 100,
		MinWorkers: 4,
		MetricInterval: func() time.Duration {
			return time.Second
		},
		MetricHandler: func(stat PoolStat) {
			log.Printf("Stat: %+s", stat)
		},
	})

	if err != nil {
		FatalFailed(t, "Error occured: %+s", err.Error())
	}

	defer wo.Shutdown()

	for i := 0; i <= 10000; i++ {
		wo.Do(fmt.Sprintf("provider: %d", i), &sim{})
	}

	LogPassed(t, "jobs were executed: total jobs")

}
