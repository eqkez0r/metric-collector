package localstorage

import (
	"context"
	"github.com/Eqke/metric-collector/pkg/metric"
	"go.uber.org/zap"
	"testing"
)

func BenchmarkLocalStorage_SetValue(b *testing.B) {
	s := New(zap.S())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SetValue(context.Background(), "gauge", "g", "234.1")
	}
}

func BenchmarkLocalStorage_SetMetric(b *testing.B) {
	s := New(zap.S())
	gauge := float64(23.41)
	m := metric.Metrics{
		ID:    "gauge",
		MType: "gauge",
		Value: &gauge,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SetMetric(context.Background(), m)
	}
}

func BenchmarkLocalStorage_SetMetrics(b *testing.B) {
	s := New(zap.S())
	gauge := float64(23.41)
	counter := int64(31)
	batch := []metric.Metrics{
		{
			ID:    "gauge",
			MType: "gauge",
			Value: &gauge,
		},
		{
			ID:    "counter",
			MType: "counter",
			Delta: &counter,
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SetMetrics(context.Background(), batch)
	}
}
