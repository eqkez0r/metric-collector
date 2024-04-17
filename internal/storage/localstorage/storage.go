package localstorage

import (
	"errors"
	store "github.com/Eqke/metric-collector/internal/storage"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/Eqke/metric-collector/pkg/metric"
	"go.uber.org/zap"
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
	logger  *zap.SugaredLogger
	mu      *sync.Mutex
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

func New(logger *zap.SugaredLogger) *LocalStorage {
	return &LocalStorage{
		storage: newStorage(),
		mu:      &sync.Mutex{},
		logger:  logger,
	}
}

func (s *LocalStorage) SetValue(metricType, name, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch metricType {
	case metric.TypeCounter.String():
		{
			metricValueInt, err := strconv.Atoi(value)
			if err != nil {
				s.logger.Error(errPointSetValue, err)
				return e.WrapError(errPointSetValue, err)
			}
			s.storage.CounterMetrics[name] += metric.Counter(metricValueInt)
		}
	case metric.TypeGauge.String():
		{
			metricGauge, err := strconv.ParseFloat(value, 64)
			if err != nil {
				s.logger.Error(errPointSetValue, err)
				return e.WrapError(errPointSetValue, err)
			}
			s.storage.GaugeMetrics[name] = metric.Gauge(metricGauge)
		}
	default:
		{
			s.logger.Error(errPointSetValue, errIsUnknownType)
			return e.WrapError(errPointSetValue, errIsUnknownType)

		}
	}
	s.logger.Infof("metric was saved with type: %s, name: %s, value: %s",
		metricType, name, value)
	return nil
}

func (s *LocalStorage) GetValue(metricType, name string) (string, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	switch metricType {
	case metric.TypeCounter.String():
		{
			if _, ok := s.storage.CounterMetrics[name]; !ok {
				s.logger.Error(errPointGetValue, errIsMetricDoesntExist)
				return "", errIsMetricDoesntExist
			}
			val := strconv.FormatInt(int64(s.storage.CounterMetrics[name]), 10)
			s.logger.Infof("metric was found with type: %s, name: %s, value: %s",
				metricType, name, val)
			return val, nil
		}
	case metric.TypeGauge.String():
		{
			if _, ok := s.storage.GaugeMetrics[name]; !ok {
				s.logger.Error(errPointGetValue, errIsMetricDoesntExist)
				return "", errIsMetricDoesntExist
			}
			val := strconv.FormatFloat(float64(s.storage.GaugeMetrics[name]), 'f', -1, 64)
			s.logger.Infof("metric was found with type: %s, name: %s, value: %s",
				metricType, name, val)
			return val, nil
		}
	default:
		{
			s.logger.Error(errPointGetValue, errIsUnknownType)
			return "", e.WrapError(errPointGetValue, errIsUnknownType)
		}
	}

}

func (s *LocalStorage) GetMetrics() ([]store.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.GaugeMetrics
}

func (s *LocalStorage) GetGaugeMetric(name string) metric.Gauge {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.GaugeMetrics[name]
}

func (s *LocalStorage) GetCounterMetrics() map[string]metric.Counter {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.CounterMetrics
}

func (s *LocalStorage) GetCounterMetric(name string) metric.Counter {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.CounterMetrics[name]
}
