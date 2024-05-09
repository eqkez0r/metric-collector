package localstorage

import (
	"bytes"
	"encoding/json"
	store "github.com/Eqke/metric-collector/internal/storage"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/Eqke/metric-collector/pkg/metric"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
)

const (
	TYPE = "Local mem database"
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
				s.logger.Error(store.ErrPointSetValue, err)
				return e.WrapError(store.ErrPointSetValue, err)
			}
			s.storage.CounterMetrics[name] += metric.Counter(metricValueInt)
		}
	case metric.TypeGauge.String():
		{
			metricGauge, err := strconv.ParseFloat(value, 64)
			if err != nil {
				s.logger.Error(store.ErrPointSetValue, err)
				return e.WrapError(store.ErrPointSetValue, err)
			}
			s.storage.GaugeMetrics[name] = metric.Gauge(metricGauge)
		}
	default:
		{
			s.logger.Error(store.ErrPointSetValue, store.ErrIsUnknownType)
			return e.WrapError(store.ErrPointSetValue, store.ErrIsUnknownType)

		}
	}
	s.logger.Infof("metric was saved with type: %s, name: %s, value: %s",
		metricType, name, value)
	return nil
}

func (s *LocalStorage) SetMetric(m metric.Metrics) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.setMetric(m)
}

func (s *LocalStorage) GetValue(metricType, name string) (string, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	switch metricType {
	case metric.TypeCounter.String():
		{
			if _, ok := s.storage.CounterMetrics[name]; !ok {
				s.logger.Error(store.ErrPointGetValue, store.ErrIsMetricDoesntExist)
				return "", store.ErrIsMetricDoesntExist
			}
			val := strconv.FormatInt(int64(s.storage.CounterMetrics[name]), 10)
			s.logger.Infof("metric was found with type: %s, name: %s, value: %s",
				metricType, name, val)
			return val, nil
		}
	case metric.TypeGauge.String():
		{
			if _, ok := s.storage.GaugeMetrics[name]; !ok {
				s.logger.Error(store.ErrPointGetValue, store.ErrIsMetricDoesntExist)
				return "", store.ErrIsMetricDoesntExist
			}
			val := strconv.FormatFloat(float64(s.storage.GaugeMetrics[name]), 'f', -1, 64)
			s.logger.Infof("metric was found with type: %s, name: %s, value: %s",
				metricType, name, val)
			return val, nil
		}
	default:
		{
			s.logger.Error(store.ErrPointGetValue, store.ErrIsUnknownType)
			return "", e.WrapError(store.ErrPointGetValue, store.ErrIsUnknownType)
		}
	}

}

func (s *LocalStorage) GetMetric(m metric.Metrics) (metric.Metrics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var met metric.Metrics
	switch m.MType {
	case metric.TypeCounter.String():
		{
			var val metric.Counter
			var ok bool
			if val, ok = s.storage.CounterMetrics[m.ID]; !ok {
				s.logger.Error(store.ErrPointGetMetric, store.ErrIsMetricDoesntExist)
				return met, store.ErrIsMetricDoesntExist
			}
			delta := int64(val)
			met = metric.Metrics{
				ID:    m.ID,
				MType: m.MType,
				Delta: &delta,
			}
		}
	case metric.TypeGauge.String():
		{
			var val metric.Gauge
			var ok bool
			if val, ok = s.storage.GaugeMetrics[m.ID]; !ok {
				s.logger.Error(store.ErrPointGetMetric, store.ErrIsMetricDoesntExist)
				return met, store.ErrIsMetricDoesntExist
			}
			value := float64(val)
			met = metric.Metrics{
				ID:    m.ID,
				MType: m.MType,
				Value: &value,
			}
		}
	default:
		{
			s.logger.Error(store.ErrPointGetMetric, store.ErrIsUnknownType)
			return met, e.WrapError(store.ErrPointGetMetric, store.ErrIsUnknownType)
		}
	}
	return met, nil
}

func (s *LocalStorage) GetMetrics() (map[string][]store.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	metrics := make(map[string][]store.Metric, 0)
	metrics[metric.TypeCounter.String()] = make([]store.Metric, 0)
	metrics[metric.TypeGauge.String()] = make([]store.Metric, 0)
	for name := range s.storage.CounterMetrics {
		m := store.Metric{
			Name:  name,
			Value: s.storage.CounterMetrics[name].String(),
		}
		metrics[metric.TypeCounter.String()] = append(metrics[metric.TypeCounter.String()], m)
	}
	for name := range s.storage.GaugeMetrics {
		m := store.Metric{
			Name:  name,
			Value: s.storage.GaugeMetrics[name].String(),
		}
		metrics[metric.TypeGauge.String()] = append(metrics[metric.TypeGauge.String()], m)
	}
	return metrics, nil
}

func (s *LocalStorage) SetMetrics(metrics []metric.Metrics) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, m := range metrics {
		err := s.setMetric(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LocalStorage) setMetric(m metric.Metrics) error {
	if m.ID == "" {
		s.logger.Error(store.ErrPointSetMetric, store.ErrIDIsEmpty)
		return store.ErrIDIsEmpty
	}

	switch m.MType {
	case metric.TypeCounter.String():
		{
			if m.Delta == nil {
				s.logger.Error(store.ErrPointSetMetric, store.ErrValueIsEmpty)
				return store.ErrValueIsEmpty
			}
			s.storage.CounterMetrics[m.ID] += metric.Counter(*m.Delta)
		}
	case metric.TypeGauge.String():
		{
			if m.Value == nil {
				s.logger.Error(store.ErrPointSetMetric, store.ErrValueIsEmpty)
				return store.ErrValueIsEmpty
			}
			s.storage.GaugeMetrics[m.ID] = metric.Gauge(*m.Value)
		}
	}
	return nil
}

func (s *LocalStorage) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s.storage, "", "  ")
}

func (s *LocalStorage) FromJSON(data []byte) error {
	return json.Unmarshal(data, &s.storage)
}

func (s *LocalStorage) ToFile(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *LocalStorage) FromFile(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var data bytes.Buffer
	_, err = f.WriteTo(&data)
	if err != nil {
		return err
	}
	return s.FromJSON(data.Bytes())
}

func (s *LocalStorage) Close() error {
	return nil
}

func (s *LocalStorage) Type() string {
	return TYPE
}
