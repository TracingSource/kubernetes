package collector

import (
	"time"

	"github.com/google/cadvisor/info/v1"
)

type FakeCollectorManager struct {
}

func (fkm *FakeCollectorManager) RegisterCollector(collector Collector) error {
	return nil
}

func (fkm *FakeCollectorManager) GetSpec() ([]v1.MetricSpec, error) {
	return []v1.MetricSpec{}, nil
}

func (fkm *FakeCollectorManager) Collect(metric map[string][]v1.MetricVal) (time.Time, map[string][]v1.MetricVal, error) {
	var zero time.Time
	return zero, metric, nil
}
