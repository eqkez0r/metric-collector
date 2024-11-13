package localstorage

import (
	"context"
	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"reflect"
	"sync"
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

func TestLocalStorage_SetValue(t *testing.T) {
	type fields struct {
		logger  *zap.SugaredLogger
		mu      *sync.Mutex
		storage storage
	}
	type args struct {
		ctx        context.Context
		metricType string
		name       string
		value      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success_set_counter",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx:        context.Background(),
				metricType: "counter",
				name:       "new_counter",
				value:      "2",
			},
			wantErr: false,
		},
		{
			name: "invalid_counter_value",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx:        context.Background(),
				metricType: "counter",
				name:       "new_counter",
				value:      "dsf.",
			},
			wantErr: true,
		},
		{
			name: "success_set_gauge",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx:        context.Background(),
				metricType: "gauge",
				name:       "new_gauge",
				value:      "2.2",
			},
			wantErr: false,
		},
		{
			name: "invalid_gauge_value",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx:        context.Background(),
				metricType: "gauge",
				name:       "new_gauge",
				value:      "dsf.",
			},
			wantErr: true,
		},
		{
			name: "invalid_type",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx:        context.Background(),
				metricType: "invalid_type",
				name:       "newcounter",
				value:      "dsf.sad",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				logger:  tt.fields.logger,
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			if err := s.SetValue(tt.args.ctx, tt.args.metricType, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalStorage_SetMetric(t *testing.T) {
	type fields struct {
		logger  *zap.SugaredLogger
		mu      *sync.Mutex
		storage storage
	}
	type args struct {
		ctx context.Context
		m   metric.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success_set_counter_metric",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "counter",
					MType: "counter",
					Delta: new(int64),
				},
			},
			wantErr: false,
		},
		{
			name: "nil_counter_value",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "counter",
					MType: "counter",
					Delta: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "success_set_gauge_metric",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "gauge",
					MType: "gauge",
					Value: new(float64),
				},
			},
			wantErr: false,
		},
		{
			name: "empty_name",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "",
					MType: "gauge",
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "nil_gauge_value",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "gauge",
					MType: "gauge",
					Value: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "unknown_metric_type",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "gauge",
					MType: "unknown",
					Value: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				logger:  tt.fields.logger,
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			if err := s.SetMetric(tt.args.ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SetMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalStorage_GetMetric(t *testing.T) {
	counter := int64(14237)
	gauge := float64(23.41)
	type fields struct {
		logger  *zap.SugaredLogger
		mu      *sync.Mutex
		storage storage
	}
	type args struct {
		ctx context.Context
		m   metric.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    metric.Metrics
		wantErr bool
	}{
		{
			name: "success_get_counter",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "counter",
					MType: "counter",
					Delta: &counter,
				},
			},
			want: metric.Metrics{
				ID:    "counter",
				MType: "counter",
				Delta: &counter,
			},
			wantErr: false,
		},
		{
			name: "success_get_gauge",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				m: metric.Metrics{
					ID:    "gauge",
					MType: "gauge",
					Value: &gauge,
				},
			},
			want: metric.Metrics{
				ID:    "gauge",
				MType: "gauge",
				Value: &gauge,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				logger:  tt.fields.logger,
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			err := s.setMetric(tt.args.ctx, tt.args.m)
			if err != nil {
				t.Errorf("setMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := s.GetMetric(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetric() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorage_GetValue(t *testing.T) {
	counter := "14"
	gauge := "323.1"
	type fields struct {
		logger  *zap.SugaredLogger
		mu      *sync.Mutex
		storage storage
	}
	type args struct {
		ctx        context.Context
		metricType string
		name       string
		value      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success_get_counter_value",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx:        context.Background(),
				metricType: "counter",
				name:       "counter",
				value:      counter,
			},
			want:    counter,
			wantErr: false,
		},
		{
			name: "success_get_gauge_value",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx:        context.Background(),
				metricType: "gauge",
				name:       "gauge",
				value:      gauge,
			},
			want:    gauge,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				logger:  tt.fields.logger,
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			err := s.SetValue(tt.args.ctx, tt.args.metricType, tt.args.name, tt.args.value)
			if err != nil {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := s.GetValue(tt.args.ctx, tt.args.metricType, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorage_GetMetrics(t *testing.T) {
	type fields struct {
		logger  *zap.SugaredLogger
		mu      *sync.Mutex
		storage storage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string][]store.Metric
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				logger:  tt.fields.logger,
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			got, err := s.GetMetrics(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if assert.NotEqual(t, nil, got) {
				t.Errorf("GetMetrics() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestLocalStorage_SetMetrics(t *testing.T) {
	counter := int64(14237)
	gauge := float64(23.41)
	type fields struct {
		logger  *zap.SugaredLogger
		mu      *sync.Mutex
		storage storage
	}
	type args struct {
		ctx     context.Context
		metrics []metric.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success_set_metrics",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				metrics: []metric.Metrics{
					{
						ID:    "counter",
						MType: "counter",
						Delta: &counter,
					},
					{
						ID:    "gauge",
						MType: "gauge",
						Value: &gauge,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failed_with_invalid_metric_type",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				metrics: []metric.Metrics{
					{
						ID:    "unknown",
						MType: "unknown",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "failed_with_empty_metrics",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				metrics: []metric.Metrics{
					{
						ID:    "",
						MType: "unknown",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "failed_with_invalid_counter",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				metrics: []metric.Metrics{
					{
						ID:    "counter",
						MType: "counter",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "failed_with_invalid_counter",
			fields: fields{
				logger:  zaptest.NewLogger(t).Sugar(),
				mu:      &sync.Mutex{},
				storage: newStorage(),
			},
			args: args{
				ctx: context.Background(),
				metrics: []metric.Metrics{
					{
						ID:    "gauge",
						MType: "gauge",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				logger:  tt.fields.logger,
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			err := s.SetMetrics(tt.args.ctx, tt.args.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
