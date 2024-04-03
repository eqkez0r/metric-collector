package metric

import (
	"math/rand"
	"runtime"
	"strconv"
)

//type Metric interface {
//	GetType() TypeMetric
//	FromString(value string)

var MetricList = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

type MetricName string

func (m MetricName) String() string {
	return string(m)
}

type TypeMetric string

func (t TypeMetric) String() string {
	return string(t)
}

const (
	Alloc         = MetricName("Alloc")
	BuckHashSys   = MetricName("BuckHashSys")
	Frees         = MetricName("Frees")
	GCCPUFraction = MetricName("GCCPUFraction")
	GCSys         = MetricName("GCSys")
	HeapAlloc     = MetricName("HeapAlloc")
	HeapIdle      = MetricName("HeapIdle")
	HeapInuse     = MetricName("HeapInuse")
	HeapObjects   = MetricName("HeapObjects")
	HeapReleased  = MetricName("HeapReleased")
	HeapSys       = MetricName("HeapSys")
	LastGC        = MetricName("LastGC")
	Lookups       = MetricName("Lookups")
	MCacheInuse   = MetricName("MCacheInuse")
	MCacheSys     = MetricName("MCacheSys")
	MSpanInuse    = MetricName("MSpanInuse")
	MSpanSys      = MetricName("MSpanSys")
	Mallocs       = MetricName("Mallocs")
	NextGC        = MetricName("NextGC")
	NumForcedGC   = MetricName("NumForcedGC")
	NumGC         = MetricName("NumGC")
	OtherSys      = MetricName("OtherSys")
	PauseTotalNs  = MetricName("PauseTotalNs")
	StackInuse    = MetricName("StackInuse")
	StackSys      = MetricName("StackSys")
	Sys           = MetricName("Sys")
	TotalAlloc    = MetricName("TotalAlloc")
	RandomValue   = MetricName("RandomValue")
	PollCount     = MetricName("PollCount")

	TypeGauge   = TypeMetric("gauge")
	TypeCounter = TypeMetric("counter")
)

func PrepareMetrics(ms *runtime.MemStats) map[TypeMetric]map[MetricName]string {
	runtime.ReadMemStats(ms)
	mp := make(map[TypeMetric]map[MetricName]string)
	mp[TypeGauge] = make(map[MetricName]string)
	mp[TypeCounter] = make(map[MetricName]string)

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

	return mp
}

func UpdateMetrics(ms *runtime.MemStats, mp map[TypeMetric]map[MetricName]string) {
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
