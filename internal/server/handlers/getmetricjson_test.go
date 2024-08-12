package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetMetricJSON(t *testing.T) {
	l := zaptest.NewLogger(t).Sugar()

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	correctGauge := 342.54
	correctCounter := int64(34)

	provider := &JSONMetricProviderMock{
		GetMetricFunc: func(contextMoqParam context.Context, metrics metric.Metrics) (metric.Metrics, error) {
			switch metrics.MType {
			case "gauge":
				{
					if metrics.ID == "random" {
						metrics.Value = &correctGauge
						return metrics, nil
					} else {
						return metric.Metrics{}, storage.ErrIsMetricDoesntExist
					}
				}
			case "counter":
				{
					if metrics.ID == "random" {
						metrics.Delta = &correctCounter
						return metrics, nil
					} else {
						return metric.Metrics{}, storage.ErrIsMetricDoesntExist
					}
				}
			default:
				{
					return metric.Metrics{}, storage.ErrIsUnknownType
				}
			}
		},
	}

	engine.GET("/value/", GetMetricJSONHandler(l, provider))

	t.Run("get_unknown_metric_json", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "some",
			MType: "unknown",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "/value/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

	})

	t.Run("get_metric_json_gauge_success", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "random",
			MType: "gauge",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "/value/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)
		m.Value = &correctGauge
		bb, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, string(bb), rr.Body.String())
	})

	t.Run("get_metric_json_gauge_error", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "unknown",
			MType: "gauge",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "/value/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

	})

	t.Run("get_metric_json_counter_success", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "random",
			MType: "counter",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "/value/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)
		m.Delta = &correctCounter
		bb, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, string(bb), rr.Body.String())
	})

	t.Run("get_metric_json_counter_error", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "unknown",
			MType: "counter",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "/value/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("get_metric_json_invalid_content_type", func(t *testing.T) {
		m := &metric.Metrics{
			ID:    "random",
			MType: "counter",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "/value/", bytes.NewReader(b))

		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

	})
}
