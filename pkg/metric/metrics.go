// Пакет metric описывает метрики и приводит их перечень
package metric

import (
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	// Перечень метрик
	PollCount = Name("PollCount")

	Alloc         = Name("Alloc")
	BuckHashSys   = Name("BuckHashSys")
	Frees         = Name("Frees")
	GCCPUFraction = Name("GCCPUFraction")
	GCSys         = Name("GCSys")
	HeapAlloc     = Name("HeapAlloc")
	HeapIdle      = Name("HeapIdle")
	HeapInuse     = Name("HeapInuse")
	HeapObjects   = Name("HeapObjects")
	HeapReleased  = Name("HeapReleased")
	HeapSys       = Name("HeapSys")
	LastGC        = Name("LastGC")
	Lookups       = Name("Lookups")
	MCacheInuse   = Name("MCacheInuse")
	MCacheSys     = Name("MCacheSys")
	MSpanInuse    = Name("MSpanInuse")
	MSpanSys      = Name("MSpanSys")
	Mallocs       = Name("Mallocs")
	NextGC        = Name("NextGC")
	NumForcedGC   = Name("NumForcedGC")
	NumGC         = Name("NumGC")
	OtherSys      = Name("OtherSys")
	PauseTotalNs  = Name("PauseTotalNs")
	StackInuse    = Name("StackInuse")
	StackSys      = Name("StackSys")
	Sys           = Name("Sys")
	TotalAlloc    = Name("TotalAlloc")
	RandomValue   = Name("RandomValue")
	TotalMemory   = Name("TotalMemory")
	FreeMemory    = Name("FreeMemory")

	// Перечень типов метрик
	TypeGauge   = MType("gauge")
	TypeCounter = MType("counter")
)

// Тип Metrics используются для API
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// Тип Name является аллиасом строки, который представляет из себя
// имя метрики
type Name string

func (m Name) String() string {
	return string(m)
}

// Тип MType является аллиасом строки, который представляет из себя
// имя тип метрики
type MType string

func (t MType) String() string {
	return string(t)
}

// Тип Map это карта карт, где первый ключ - это тип метрики,
// а второй ключ - имя метрики.
type Map map[MType]map[Name]string

// Функция UpdateRuntimeMetrics выполянет обновление метрик.
// ms - переменная, которая хранит MemStats
// mp - карта, куда записываются метрики
func UpdateRuntimeMetrics(ms *runtime.MemStats, mp Map) {
	runtime.ReadMemStats(ms)
	mp[TypeGauge][Alloc] = strconv.FormatFloat(float64(ms.Alloc), 'f', -1, 64)
	mp[TypeGauge][BuckHashSys] = strconv.FormatFloat(float64(ms.BuckHashSys), 'f', -1, 64)
	mp[TypeGauge][Frees] = strconv.FormatFloat(float64(ms.Frees), 'f', -1, 64)
	mp[TypeGauge][GCCPUFraction] = strconv.FormatFloat(float64(ms.GCCPUFraction), 'f', -1, 64)
	mp[TypeGauge][GCSys] = strconv.FormatFloat(float64(ms.GCSys), 'f', -1, 64)
	mp[TypeGauge][HeapAlloc] = strconv.FormatFloat(float64(ms.HeapAlloc), 'f', -1, 64)
	mp[TypeGauge][HeapIdle] = strconv.FormatFloat(float64(ms.HeapIdle), 'f', -1, 64)
	mp[TypeGauge][HeapInuse] = strconv.FormatFloat(float64(ms.HeapInuse), 'f', -1, 64)
	mp[TypeGauge][HeapObjects] = strconv.FormatFloat(float64(ms.HeapObjects), 'f', -1, 64)
	mp[TypeGauge][HeapReleased] = strconv.FormatFloat(float64(ms.HeapReleased), 'f', -1, 64)
	mp[TypeGauge][HeapSys] = strconv.FormatFloat(float64(ms.HeapSys), 'f', -1, 64)
	mp[TypeGauge][LastGC] = strconv.FormatFloat(float64(ms.LastGC), 'f', -1, 64)
	mp[TypeGauge][Lookups] = strconv.FormatFloat(float64(ms.Lookups), 'f', -1, 64)
	mp[TypeGauge][MCacheInuse] = strconv.FormatFloat(float64(ms.MCacheInuse), 'f', -1, 64)
	mp[TypeGauge][MCacheSys] = strconv.FormatFloat(float64(ms.MCacheSys), 'f', -1, 64)
	mp[TypeGauge][MSpanInuse] = strconv.FormatFloat(float64(ms.MSpanInuse), 'f', -1, 64)
	mp[TypeGauge][MSpanSys] = strconv.FormatFloat(float64(ms.MSpanSys), 'f', -1, 64)
	mp[TypeGauge][Mallocs] = strconv.FormatFloat(float64(ms.Mallocs), 'f', -1, 64)
	mp[TypeGauge][NextGC] = strconv.FormatFloat(float64(ms.NextGC), 'f', -1, 64)
	mp[TypeGauge][NumForcedGC] = strconv.FormatFloat(float64(ms.NumForcedGC), 'f', -1, 64)
	mp[TypeGauge][NumGC] = strconv.FormatFloat(float64(ms.NumGC), 'f', -1, 64)
	mp[TypeGauge][OtherSys] = strconv.FormatFloat(float64(ms.OtherSys), 'f', -1, 64)
	mp[TypeGauge][PauseTotalNs] = strconv.FormatFloat(float64(ms.PauseTotalNs), 'f', -1, 64)
	mp[TypeGauge][StackInuse] = strconv.FormatFloat(float64(ms.StackInuse), 'f', -1, 64)
	mp[TypeGauge][StackSys] = strconv.FormatFloat(float64(ms.StackSys), 'f', -1, 64)
	mp[TypeGauge][Sys] = strconv.FormatFloat(float64(ms.Sys), 'f', -1, 64)
	mp[TypeGauge][TotalAlloc] = strconv.FormatFloat(float64(ms.TotalAlloc), 'f', -1, 64)
	mp[TypeGauge][RandomValue] = strconv.FormatFloat(rand.Float64(), 'f', -1, 64)

}

// Функция UpdateUtilMetrics позволяет получить дополнительные метрики.
// mp - карта, куда будут записаны метрики.
func UpdateUtilMetrics(mp Map) error {
	m, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	cpu, err := cpu.Percent(time.Second, true)
	if err != nil {
		return err
	}
	mp[TypeGauge][TotalMemory] = strconv.FormatUint(m.Free, 10)
	mp[TypeGauge][FreeMemory] = strconv.FormatUint(m.Free, 10)
	for i, v := range cpu {
		name := Name("CPUutilization" + strconv.Itoa(i))
		mp[TypeGauge][name] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return nil
}
