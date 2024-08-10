package handlers

import (
	"context"
	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPostMetric(t *testing.T) {
	l := zaptest.NewLogger(t).Sugar()

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	provider := &NewMetricProviderMock{
		SetValueFunc: func(contextMoqParam context.Context, mtype, name, value string) error {
			if name == "" {
				return store.ErrIDIsEmpty
			}
			switch mtype {
			case "gauge":
				{
					_, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return err
					}
					return nil
				}
			case "counter":
				{
					_, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return err
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

	engine.POST("/update/:type/:name/:value/", POSTMetricHandler(l, provider))

	t.Run("post_counter_success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update/counter/random/1/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("post_counter_error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update/counter/random/fsd/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("post_gauge_success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update/gauge/random/234.35/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("post_gauge_error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update/counter/random/fsd/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("post_unknown_type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update/unknown/random/234.35/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("post_with_empty_name", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update/gauge//1/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

}
