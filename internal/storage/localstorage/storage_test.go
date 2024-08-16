package localstorage

import (
	"context"
	"go.uber.org/zap/zaptest"
	"testing"

	"github.com/Eqke/metric-collector/pkg/metric"
	"go.uber.org/zap"
)

func BenchmarkLocalStorage_SetValue(b *testing.B) {
	l := zaptest.NewLogger(b).Sugar()
	s := New(l)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := s.SetValue(context.Background(), "gauge", "g", "234.1")
		if err != nil {
			l.Error("Error setting value", zap.Error(err))
		}
	}
}

func BenchmarkLocalStorage_SetMetric(b *testing.B) {
	l := zaptest.NewLogger(b).Sugar()
	s := New(l)
	gauge := float64(23.41)
	m := metric.Metrics{
		ID:    "gauge",
		MType: "gauge",
		Value: &gauge,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := s.SetMetric(context.Background(), m)
		if err != nil {
			l.Error("Error setting metric", zap.Error(err))
		}
	}
}

func BenchmarkLocalStorage_SetMetrics(b *testing.B) {
	l := zaptest.NewLogger(b).Sugar()
	s := New(l)
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
		err := s.SetMetrics(context.Background(), batch)
		if err != nil {
			l.Error("Error setting metrics", zap.Error(err))
		}
	}
}
