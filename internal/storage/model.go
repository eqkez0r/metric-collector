// Пакет storage предоставляет интерфейс для хранилища.
package storage

import (
	"context"
	"github.com/Eqke/metric-collector/pkg/metric"
)

// Интерфейс хранилища
type Storage interface {
	// Метод SetValue добавляет или обновляет значение метрики.
	// Получает на вход(порядок соответствует):
	// metricType - тип метрики,
	// name - имя метрики,
	// value - значение метрики.
	SetValue(context.Context, string, string, string) error

	// Метод SetMetric добавляет или обновляет значение метрики с
	// помощью типа metric.Metrics.
	// Получает экземпляр metric.Metrics.
	SetMetric(context.Context, metric.Metrics) error

	// Метод SetMetrics добавляет или обновляет массив метрик
	// Получает на вход массив metric.Metrics.
	SetMetrics(context.Context, []metric.Metrics) error

	// Метод GetValue позволяет получить значение метрики.
	// Получает на вход:
	// metricType - тип метрики,
	// name - имя метрики.
	GetValue(context.Context, string, string) (string, error)

	// Метод GetMetric полвзоялет получить экземлпяр metric.Metric.
	// Получает на вход:
	// m - экземпляр metric.Metrics
	GetMetric(context.Context, metric.Metrics) (metric.Metrics, error)

	// Метод GetMetrics позволяет получить карту метрик
	GetMetrics(context.Context) (map[string][]Metric, error)

	// Метод ToJSON используется для сериализации
	ToJSON(context.Context) ([]byte, error)

	// Метод FromJSON используется для десериализации
	FromJSON(context.Context, []byte) error

	// Метод ToFile используется для записи в файл
	ToFile(context.Context, string) error

	// Метод FromFile используется для чтения из файла
	FromFile(context.Context, string) error

	// Метод Type используется для получения типа хранилища
	Type() string

	// Метод Close используется для утилизации ресурсов
	Close() error
}

// Тип Metric нужен используется для дальнейшего отображения
// в строковом формате
type Metric struct {
	Name  string
	Value string
}
