package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	errPointPostMetric = "error in POST /update/:type/:name/:value"
)

func POSTMetricHandler(
	logger *zap.SugaredLogger,
	storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		metricType := context.Param("type")
		if metricType != metric.TypeCounter.String() && metricType != metric.TypeGauge.String() {
			logger.Errorf("%s: unknown metric type %s", errPointPostMetric, metricType)
			context.Status(http.StatusBadRequest)
			return
		}
		metricName := context.Param("name")
		if metricName == "" {
			logger.Errorf("%s: empty metric name", errPointPostMetric)
			context.Status(http.StatusNotFound)
			return
		}
		metricValue := context.Param("value")
		logger.Infof("metric was received with type: %s, name: %s, value: %s", metricType, metricName, metricValue)
		err := storage.SetValue(metricType, metricName, metricValue)
		if err != nil {
			logger.Errorf("%s: %v", errPointPostMetric, err)
			context.Status(http.StatusBadRequest)
			return
		}
		logger.Infof("metric was saved with type: %s, name: %s, value: %s", metricType, metricName, metricValue)
		context.Status(http.StatusOK)

	}
}
