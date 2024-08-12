package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	errPointGetRootMetrics = "error in GET /"
)

//go:generate moq -out rootMetricProvider_moq_test.go . RootMetricsProvider
type RootMetricsProvider interface {
	GetMetrics(context.Context) (map[string][]store.Metric, error)
}

func GetRootMetricsHandler(
	logger *zap.SugaredLogger,
	storage RootMetricsProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("root metrics handler was called")
		var data map[string][]store.Metric
		var err error
		err = retry.Retry(logger, 3, func() error {
			data, err = storage.GetMetrics(c)
			return err
		})
		if err != nil {
			logger.Errorf("%s: %v", errPointGetRootMetrics, err)
			c.Status(http.StatusInternalServerError)
			return
		}
		bytes, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			logger.Errorf("%s: %v", errPointGetRootMetrics, err)
			c.Status(http.StatusInternalServerError)
			return
		}
		logger.Info("get metrics was success")

		c.Data(http.StatusOK, "html/text", bytes)
	}
}
