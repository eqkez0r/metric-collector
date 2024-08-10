package handlers

import (
	"context"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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
	return func(context *gin.Context) {
		logger.Infof("/update/:type/:name/:value post metric")
		metricType := context.Param("type")
		metricName := context.Param("name")

		if metricType != metric.TypeCounter.String() && metricType != metric.TypeGauge.String() {
			logger.Errorf("%s: unknown metric type %s", errPointPostMetric, metricType)
			context.Status(http.StatusBadRequest)
			return
		}

		if metricName == "" {
			logger.Errorf("%s: empty metric name", errPointPostMetric)
			context.Status(http.StatusNotFound)
			return
		}
		metricValue := context.Param("value")
		logger.Infof("metric was received with type: %s, name: %s, value: %s",
			metricType, metricName, metricValue)
		err := retry.Retry(logger, 3, func() error {
			return s.SetValue(context, metricType, metricName, metricValue)
		})
		if err != nil {
			logger.Errorf("%s: %v", errPointPostMetric, err)
			context.Status(http.StatusBadRequest)
			return
		}
		logger.Infof("metric was saved with type: %s, name: %s, value: %s",
			metricType, metricName, metricValue)

		context.Status(http.StatusOK)

	}
}
