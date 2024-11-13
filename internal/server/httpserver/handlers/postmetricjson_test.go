package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestPostMetricJSON(t *testing.T) {
	l := zaptest.NewLogger(t).Sugar()

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	provider := &NewJSONMetricProviderMock{
		SetMetricFunc: func(contextMoqParam context.Context, metric metric.Metrics) error {
			if metric.ID == "" {
				return store.ErrIDIsEmpty
			}
			switch metric.MType {
			case "gauge":
				{
					if metric.Value == nil {
						return store.ErrValueIsEmpty
					}
					return nil
				}
			case "counter":
				{
					if metric.Delta == nil {
						return store.ErrValueIsEmpty
					}
					return nil
				}
			default:
				{
					return store.ErrIsUnknownType
				}
			}
		},
	}

	engine.POST("/update/", POSTMetricJSONHandler(l, provider))

	t.Run("post_gauge_success", func(t *testing.T) {
		val := float64(3414.2)
		m := &metric.Metrics{
			ID:    "some",
			MType: "gauge",
			Value: &val,
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/update/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("post_gauge_error", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "some",
			MType: "gauge",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/update/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("post_counter_success", func(t *testing.T) {
		delta := int64(3414)
		m := &metric.Metrics{
			ID:    "some",
			MType: "counter",
			Delta: &delta,
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/update/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("post_counter_error", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "some",
			MType: "counter",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/update/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("post_with_empty_name", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "",
			MType: "gauge",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/update/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("post_unknown_type", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "some",
			MType: "unknown",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/update/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("post_invalid_json", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/update/", bytes.NewReader([]byte("{")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

//func TestPOSTMetricJSONHandler(t *testing.T) {
//	type want struct {
//		statusCode int
//	}
//	delta := int64(10)
//	value := float64(210.3)
//	tests := []struct {
//		name string
//		url  string
//		m    metric.Metrics
//		want want
//	}{
//		{
//			name: "Success request with counter",
//			url:  "http://localhost:8080/update",
//			m:    metric.Metrics{ID: "test", MType: "counter", Delta: &delta},
//			want: want{
//				statusCode: http.StatusOK,
//			},
//		},
//		{
//			name: "Success request with gauger",
//			url:  "http://localhost:8080/update",
//			m:    metric.Metrics{ID: "test", MType: "gauge", Value: &value},
//			want: want{
//				statusCode: http.StatusOK,
//			},
//		},
//		{
//			name: "Without metric name",
//			url:  "http://localhost:8080/update",
//			m:    metric.Metrics{ID: "", MType: "gauge", Value: &value},
//			want: want{
//				statusCode: http.StatusNotFound,
//			},
//		},
//		{
//			name: "Without counter value",
//			url:  "http://localhost:8080/update",
//			m:    metric.Metrics{ID: "test", MType: "counter", Delta: nil, Value: nil},
//			want: want{
//				statusCode: http.StatusNotFound,
//			},
//		},
//		{
//			name: "Without gauge value",
//			url:  "http://localhost:8080/update",
//			m:    metric.Metrics{ID: "test", MType: "gauge", Delta: nil, Value: nil},
//			want: want{
//				statusCode: http.StatusNotFound,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			t.Run(tt.name, func(t *testing.T) {
//				l := zaptest.NewLogger(t).Sugar()
//				localstorage := ls.New(l)
//				val, _ := json.Marshal(tt.m)
//				req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader(val))
//				req.Header.Set("Content-Type", "application/json")
//				w := httptest.NewRecorder()
//				gin.DisableConsoleColor()
//				gin.SetMode(gin.ReleaseMode)
//
//				engine := gin.New()
//				engine.POST("/update", POSTMetricJSONHandler(l, localstorage))
//				engine.ServeHTTP(w, req)
//
//				res := w.Result()
//				defer res.Body.Close()
//
//				assert.Equal(t, tt.want.statusCode, res.StatusCode)
//			})
//		})
//	}
//}
