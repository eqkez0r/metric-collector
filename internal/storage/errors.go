// Пакет storage предоставляет интерфейс для хранилища.
package storage

import "errors"

// Перечень ошибок и точек для ошибок
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
