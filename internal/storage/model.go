package storage

type Metric struct {
	Name  string
	Value string
}

type Storage interface {
	SetValue(metricType, name, value string) error
	GetValue(metricType, name string) (string, error)
	GetMetrics() ([]Metric, error)
}
