package localstorage

import (
	"errors"
	store "github.com/Eqke/metric-collector/internal/storage"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/Eqke/metric-collector/pkg/metric"
	"log"
	"strconv"
	"sync"
)

var (
	errIsMetricDoesntExist = errors.New("metric doesn't exist")
	errIsUnknownType       = errors.New("unknown metric type")

	errPointSetValue = "error in localstorage.SetValue(): "
	errPointGetValue = "error in localstorage.GetValue(): "
)

type LocalStorage struct {
	mu      *sync.RWMutex
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
		mu:      &sync.RWMutex{},
	}
}

func (s *LocalStorage) SetValue(metricType, name, value string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch metricType {
	case metric.TypeCounter.String():
		{
			metricValueInt, err := strconv.Atoi(value)
			if err != nil {
				log.Println(errPointSetValue, err)
				return e.WrapError(errPointSetValue, err)
			}
			s.storage.CounterMetrics[name] += metric.Counter(metricValueInt)
		}
	case metric.TypeGauge.String():
		{
			metricGauge, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Println(errPointSetValue, err)
				return e.WrapError(errPointSetValue, err)
			}
			s.storage.GaugeMetrics[name] = metric.Gauge(metricGauge)
		}
	default:
		{
			log.Println(errPointSetValue, errIsUnknownType)
			return e.WrapError(errPointSetValue, errIsUnknownType)

		}
	}
	return nil
}

func (s *LocalStorage) GetValue(metricType, name string) (string, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	switch metricType {
	case metric.TypeCounter.String():
		{
			if _, ok := s.storage.CounterMetrics[name]; !ok {
				log.Println(errPointGetValue, errIsMetricDoesntExist)
				return "", errIsMetricDoesntExist
			}
			val := strconv.FormatInt(int64(s.storage.CounterMetrics[name]), 10)
			log.Println("metric was found", metricType, name, val)

			return val, nil
		}
	case metric.TypeGauge.String():
		{
			if _, ok := s.storage.GaugeMetrics[name]; !ok {
				log.Println(errPointGetValue, errIsMetricDoesntExist)
				return "", errIsMetricDoesntExist
			}
			val := strconv.FormatFloat(float64(s.storage.GaugeMetrics[name]), 'f', -1, 64)
			log.Println("metric was found", metricType, name, val)
			return val, nil
		}
	}
	return "", nil
}

func (s *LocalStorage) GetMetrics() ([]store.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	metrics := make([]store.Metric, 0, len(s.storage.CounterMetrics)+len(s.storage.GaugeMetrics))
	for name := range s.storage.CounterMetrics {
		m := store.Metric{
			Name:  name,
			Value: s.storage.CounterMetrics[name].String(),
		}
		metrics = append(metrics, m)
	}
	for name := range s.storage.GaugeMetrics {
		m := store.Metric{
			Name:  name,
			Value: s.storage.GaugeMetrics[name].String(),
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (s *LocalStorage) GetGaugeMetrics() map[string]metric.Gauge {
	return s.storage.GaugeMetrics
}

func (s *LocalStorage) GetGaugeMetric(name string) metric.Gauge {
	return s.storage.GaugeMetrics[name]
}

func (s *LocalStorage) GetCounterMetrics() map[string]metric.Counter {
	return s.storage.CounterMetrics
}

func (s *LocalStorage) GetCounterMetric(name string) metric.Counter {
	return s.storage.CounterMetrics[name]
}
