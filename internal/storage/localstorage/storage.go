package localstorage

import (
	"github.com/Eqke/metric-collector/pkg/metric"
	"log"
	"sync"
)

type LocalStorage struct {
	mu      sync.RWMutex
	storage storage
}

type storage struct {
	// <NameMetric, Metric>
	GaugeMetrics   map[string]metric.Gauge
	CounterMetrics map[string]metric.Counter
}

func newStorage() storage {
	//share for new metric
	return storage{
		GaugeMetrics:   make(map[string]metric.Gauge),
		CounterMetrics: make(map[string]metric.Counter),
	}
}

func New() *LocalStorage {
	return &LocalStorage{
		storage: newStorage(),
	}
}

func (s *LocalStorage) Gauge(name string, value metric.Gauge) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage.GaugeMetrics[name] = value
	log.Printf("Gauge: %s, value: %f", name, s.storage.GaugeMetrics[name])
	return nil
}

func (s *LocalStorage) Counter(name string, value metric.Counter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage.CounterMetrics[name] += value
	log.Printf("Counter: %s, value: %d", name, s.storage.CounterMetrics[name])
	return nil
}
