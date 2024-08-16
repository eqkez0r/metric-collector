// Пакет metric описывает метрики и приводит их перечень
package metric

import "strconv"

// Тип Gauge является аллисасом для float64
type Gauge float64

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}
