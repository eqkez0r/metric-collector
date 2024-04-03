package storage

type Storage interface {
	Gauge(name string, value float64) error
	Counter(name string, value int64) error
}
