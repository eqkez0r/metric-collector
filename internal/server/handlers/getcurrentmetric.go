package handlers

import (
	"context"
	"net/http"

	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	errPointGetMetric = "error in GET /value/:type/:name"
)

//go:generate moq -out currentMetricProvider_moq_test.go . CurrentMetricProvider
type CurrentMetricProvider interface {
	GetValue(context.Context, string, string) (string, error)
}

func GETMetricHandler(
	logger *zap.SugaredLogger,
	s CurrentMetricProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("/value/:type/:name get metric")
		metricType := c.Param("type")
		metricName := c.Param("name")
		logger.Infof("metric was requested with type: %s, name: %s", metricType, metricName)
		var value string
		var err error
		err = retry.Retry(logger, 3, func() error {
			value, err = s.GetValue(c, metricType, metricName)
			return err
		})
		if err != nil {
			logger.Errorf("%s: %v", errPointGetMetric, err)
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, value)
		logger.Infof("metric get success with type: %s, name: %s, value: %s", metricType, metricName, value)
	}
}
