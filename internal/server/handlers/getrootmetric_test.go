package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestRootMetric(t *testing.T) {
	t.Run("get_root_metric", func(t *testing.T) {

		l := zaptest.NewLogger(t).Sugar()

		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
		engine := gin.New()
		engine.RedirectFixedPath = true

		metrics := map[string][]store.Metric{
			"gauge": {
				{
					Name:  "some_gauge",
					Value: "342.42",
				},
			},
			"counter": {
				{
					Name:  "some_counter",
					Value: "34543",
				},
			},
		}

		provider := &RootMetricsProviderMock{
			GetMetricsFunc: func(contextMoqParam context.Context) (map[string][]store.Metric, error) {
				return metrics, nil
			},
		}

		engine.GET("/", GetRootMetricsHandler(l, provider))

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		rr := httptest.NewRecorder()
		require.Equal(t, http.StatusOK, rr.Code)

	})
}
