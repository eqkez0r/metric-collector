package handlers

import (
	"context"
	"net/http"

	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	errPointPostMetric = "error in POST /update/:type/:name/:value"
)

//go:generate moq -out newMetricProvider_moq_test.go . NewMetricProvider
type NewMetricProvider interface {
	SetValue(context.Context, string, string, string) error
}

func POSTMetricHandler(
	logger *zap.SugaredLogger,
	s NewMetricProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("/update/:type/:name/:value post metric")
		metricType := c.Param("type")
		metricName := c.Param("name")

		if metricType != metric.TypeCounter.String() && metricType != metric.TypeGauge.String() {
			logger.Errorf("%s: unknown metric type %s", errPointPostMetric, metricType)
			c.Status(http.StatusBadRequest)
			return
		}

		if metricName == "" {
			logger.Errorf("%s: empty metric name", errPointPostMetric)
			c.Status(http.StatusNotFound)
			return
		}
		metricValue := c.Param("value")
		logger.Infof("metric was received with type: %s, name: %s, value: %s",
			metricType, metricName, metricValue)
		err := retry.Retry(logger, 3, func() error {
			return s.SetValue(c, metricType, metricName, metricValue)
		})
		if err != nil {
			logger.Errorf("%s: %v", errPointPostMetric, err)
			c.Status(http.StatusBadRequest)
			return
		}
		logger.Infof("metric was saved with type: %s, name: %s, value: %s",
			metricType, metricName, metricValue)

		c.Status(http.StatusOK)

	}
}
