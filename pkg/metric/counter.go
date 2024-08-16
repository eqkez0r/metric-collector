// Пакет metric описывает метрики и приводит их перечень
package metric

import "strconv"

// Тип Counter представляет аллиас на int64
type Counter int64

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}
