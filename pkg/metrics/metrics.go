package metrics

//type Metric interface {
//	GetType() TypeMetric
//	FromString(value string)
//}

type TypeMetric string

const (
	TypeMetricsCounter TypeMetric = "counter"
	TypeMetricsGauge   TypeMetric = "gauge"
)
