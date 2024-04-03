package localstorage

import (
	"errors"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"log"
	"strconv"
	"sync"
)

var (
	errIsMetricDoesntExist = errors.New("metric doesn't exist")
	errIsUnknownType       = errors.New("unknown metric type")
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

func (s *LocalStorage) SetValue(metricType, name, value string) error {
	const errPoint = "error in localstorage.SetValue(): "
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch metricType {
	case metric.TypeCounter.String():
		{
			metricValueInt, err := strconv.Atoi(value)
			if err != nil {
				log.Println(errPoint, err)
				return err
			}
			s.storage.CounterMetrics[name] += metric.Counter(metricValueInt)
		}
	case metric.TypeGauge.String():
		{
			metricGauge, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Println(errPoint, err)
				return err
			}
			s.storage.GaugeMetrics[name] = metric.Gauge(metricGauge)
		}
	default:
		{
			log.Println(errPoint, errIsUnknownType)
			return errIsUnknownType

		}
	}
	return nil
}

func (s *LocalStorage) GetValue(metricType, name string) (string, error) {
	const errPoint = "error in localstorage.GetValue(): "
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch metricType {
	case metric.TypeCounter.String():
		{
			if _, ok := s.storage.CounterMetrics[name]; !ok {
				log.Println(errPoint, errIsMetricDoesntExist)
				return "", errIsMetricDoesntExist
			}
			val := strconv.FormatInt(int64(s.storage.CounterMetrics[name]), 10)
			log.Println("metric was found", metricType, name, val)

			return val, nil
		}
	case metric.TypeGauge.String():
		{
			if _, ok := s.storage.GaugeMetrics[name]; !ok {
				log.Println(errPoint, errIsMetricDoesntExist)
				return "", errIsMetricDoesntExist
			}
			val := strconv.FormatFloat(float64(s.storage.GaugeMetrics[name]), 'f', -1, 64)
			log.Println("metric was found", metricType, name, val)
			return val, nil
		}
	}
	return "", nil
}

func (s *LocalStorage) GetMetrics() ([]stor.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	metrics := make([]stor.Metric, 0, len(s.storage.CounterMetrics)+len(s.storage.GaugeMetrics))
	for name := range s.storage.CounterMetrics {
		m := stor.Metric{
			Name:  name,
			Value: s.storage.CounterMetrics[name].String(),
		}
		metrics = append(metrics, m)
	}
	for name := range s.storage.GaugeMetrics {
		m := stor.Metric{
			Name:  name,
			Value: s.storage.GaugeMetrics[name].String(),
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}
