package storage

import (
	"errors"
	"github.com/Eqke/metric-collector/pkg/metric"
)

var (
	ErrIsMetricDoesntExist = errors.New("metric doesn't exist")
	ErrIsUnknownType       = errors.New("unknown metric type")
	ErrIDIsEmpty           = errors.New("metric name is empty")
	ErrValueIsEmpty        = errors.New("metric value is empty")
)

type Metric struct {
	Name  string
	Value string
}

type Storage interface {
	SetValue(string, string, string) error
	SetMetric(metric.Metrics) error
	GetValue(string, string) (string, error)
	GetMetrics() ([]Metric, error)
	GetMetric(metric.Metrics) (metric.Metrics, error)
	GetGaugeMetrics() map[string]metric.Gauge
	GetGaugeMetric(string) metric.Gauge
	GetCounterMetrics() map[string]metric.Counter
	GetCounterMetric(string) metric.Counter
	ToJSON() ([]byte, error)
	FromJSON([]byte) error
	ToFile(string) error
	FromFile(string) error
}
