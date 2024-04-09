package storage

import "github.com/Eqke/metric-collector/pkg/metric"

type Metric struct {
	Name  string
	Value string
}

type Storage interface {
	SetValue(metricType, name, value string) error
	GetValue(metricType, name string) (string, error)
	GetMetrics() ([]Metric, error)
	GetGaugeMetrics() map[string]metric.Gauge
	GetGaugeMetric(name string) metric.Gauge
	GetCounterMetrics() map[string]metric.Counter
	GetCounterMetric(name string) metric.Counter
}
