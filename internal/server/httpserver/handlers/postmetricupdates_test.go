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

func TestPostMetricUpdates(t *testing.T) {
	l := zaptest.NewLogger(t).Sugar()

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	provider := &BatchMetricProviderMock{
		SetMetricsFunc: func(contextMoqParam context.Context, ms []metric.Metrics) error {
			for _, m := range ms {
				if m.ID == "" {
					return store.ErrIDIsEmpty
				}
				switch m.MType {
				case "gauge":
					{
						if m.Value == nil {
							return store.ErrValueIsEmpty
						}
					}
				case "counter":
					{
						if m.Delta == nil {
							return store.ErrValueIsEmpty
						}
					}
				default:
					{
						return store.ErrIsUnknownType
					}
				}
			}
			return nil
		},
	}

	gauge := float64(213.3)
	counter := int64(84)

	engine.POST("/updates/", PostMetricUpdates(l, provider))

	t.Run("update_success", func(t *testing.T) {

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

		b, err := json.Marshal(batch)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("update_error", func(t *testing.T) {

		t.Run("bad_content_type", func(t *testing.T) {

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

			b, err := json.Marshal(batch)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("gauge_error", func(t *testing.T) {
			batch := []metric.Metrics{
				{
					ID:    "gauge",
					MType: "gauge",
				},
				{
					ID:    "counter",
					MType: "counter",
					Delta: &counter,
				},
			}
			b, err := json.Marshal(batch)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("counter_error", func(t *testing.T) {
			batch := []metric.Metrics{
				{
					ID:    "gauge",
					MType: "gauge",
					Value: &gauge,
				},
				{
					ID:    "counter",
					MType: "counter",
				},
			}
			b, err := json.Marshal(batch)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("empty_name", func(t *testing.T) {
			batch := []metric.Metrics{
				{
					ID:    "gauge",
					MType: "gauge",
					Value: &gauge,
				},
				{
					ID:    "",
					MType: "counter",
					Delta: &counter,
				},
			}
			b, err := json.Marshal(batch)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("unknown_type", func(t *testing.T) {
			batch := []metric.Metrics{
				{
					ID:    "gauge",
					MType: "unknown",
					Value: &gauge,
				},
				{
					ID:    "counter",
					MType: "counter",
					Delta: &counter,
				},
			}
			b, err := json.Marshal(batch)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	})

}
