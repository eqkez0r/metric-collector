package handlers

import (
	"bytes"
	"encoding/json"
	ls "github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPOSTMetricJSONHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	delta := int64(10)
	value := float64(210.3)
	tests := []struct {
		name string
		url  string
		m    metric.Metrics
		want want
	}{
		{
			name: "Success request with counter",
			url:  "http://localhost:8080/update",
			m:    metric.Metrics{ID: "test", MType: "counter", Delta: &delta},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Success request with gauger",
			url:  "http://localhost:8080/update",
			m:    metric.Metrics{ID: "test", MType: "gauge", Value: &value},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Without metric name",
			url:  "http://localhost:8080/update",
			m:    metric.Metrics{ID: "", MType: "gauge", Value: &value},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Without counter value",
			url:  "http://localhost:8080/update",
			m:    metric.Metrics{ID: "test", MType: "counter", Delta: nil, Value: nil},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Without gauge value",
			url:  "http://localhost:8080/update",
			m:    metric.Metrics{ID: "test", MType: "gauge", Delta: nil, Value: nil},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				l := zaptest.NewLogger(t).Sugar()
				localstorage := ls.New(l)
				val, _ := json.Marshal(tt.m)
				req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader(val))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				gin.DisableConsoleColor()
				gin.SetMode(gin.ReleaseMode)

				engine := gin.New()
				engine.POST("/update", POSTMetricJSONHandler(l, localstorage))
				engine.ServeHTTP(w, req)

				res := w.Result()
				defer res.Body.Close()

				assert.Equal(t, tt.want.statusCode, res.StatusCode)
			})
		})
	}
}
