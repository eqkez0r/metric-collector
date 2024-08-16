package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetCurrentMetric(t *testing.T) {

	l := zaptest.NewLogger(t).Sugar()

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	provider := &CurrentMetricProviderMock{
		GetValueFunc: func(contextMoqParam context.Context, typeM, name string) (string, error) {
			switch typeM {
			case "gauge":
				{
					if name == "random" {
						return "342.45", nil
					} else {
						return "", storage.ErrIsMetricDoesntExist
					}
				}
			case "counter":
				{
					if name == "random" {
						return "34", nil
					} else {
						return "", storage.ErrIsMetricDoesntExist
					}
				}
			default:
				{
					return "", storage.ErrIsUnknownType
				}
			}
		},
	}
	engine.GET("/value/:type/:name/", GETMetricHandler(l, provider))

	t.Run("get_unknown_metric", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/value/sdsad/unknown/", nil)
		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("get_gauge_metric_success", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/value/gauge/random/", nil)
		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "342.45", rr.Body.String())

	})

	t.Run("get_gauge_metric_error", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/value/gauge/notexist/", nil)
		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

	})

	t.Run("get_counter_metric_success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/value/counter/random/", nil)
		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "34", rr.Body.String())
	})

	t.Run("get_counter_metric_error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/value/counter/notexist/", nil)
		rr := httptest.NewRecorder()

		engine.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)
	})
}
