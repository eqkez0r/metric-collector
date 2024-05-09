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

	ErrPointSetValue          = "error in storage.SetValue(): "
	ErrPointSetMetric         = "error in storage.SetMetric(): "
	ErrPointGetValue          = "error in storage.GetValue(): "
	ErrPointGetMetrics        = "error in storage.GetMetrics(): "
	ErrPointGetMetric         = "error in storage.GetMetric(): "
	ErrPointGetGaugeMetrics   = "error in storage.GetGaugeMetrics(): "
	ErrPointGetGaugeMetric    = "error in storage.GetGaugeMetric(): "
	ErrPointGetCounterMetrics = "error in storage.GetCounterMetrics(): "
	ErrPointGetCounterMetric  = "error in storage.GetCounterMetric(): "
	ErrPointToJSON            = "error in storage.ToJSON(): "
	ErrPointFromJSON          = "error in storage.FromJSON(): "
	ErrPointToFile            = "error in storage.ToFile(): "
	ErrPointFromFile          = "error in storage.FromFile(): "
	ErrPointClose             = "error in storage.Close(): "
)

type Metric struct {
	Name  string
	Value string
}

type Storage interface {
	SetValue(string, string, string) error
	SetMetric(metric.Metrics) error
	GetValue(string, string) (string, error)
	GetMetrics() (map[string][]Metric, error)
	GetMetric(metric.Metrics) (metric.Metrics, error)
	SetMetrics([]metric.Metrics) error
	ToJSON() ([]byte, error)
	FromJSON([]byte) error
	ToFile(string) error
	FromFile(string) error
	Type() string
	Close() error
}
