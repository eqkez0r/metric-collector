package storage

import (
	"context"
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
	SetValue(context.Context, string, string, string) error
	SetMetric(context.Context, metric.Metrics) error
	GetValue(context.Context, string, string) (string, error)
	GetMetrics(context.Context) (map[string][]Metric, error)
	GetMetric(context.Context, metric.Metrics) (metric.Metrics, error)
	SetMetrics(context.Context, []metric.Metrics) error
	ToJSON(context.Context) ([]byte, error)
	FromJSON(context.Context, []byte) error
	ToFile(context.Context, string) error
	FromFile(context.Context, string) error
	Type() string
	Close() error
}
