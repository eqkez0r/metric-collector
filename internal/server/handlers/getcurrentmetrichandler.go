package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	errPointGetMetric = "error in GET /value/:type/:name"
)

func GETMetricHandler(
	logger *zap.SugaredLogger,
	storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Debug("get metric handler was called")
		metricType := context.Param("type")
		metricName := context.Param("name")
		logger.Infof("metric was requested with type: %s, name: %s", metricType, metricName)
		value, err := storage.GetValue(metricType, metricName)
		if err != nil {
			logger.Errorf("%s: %v", errPointGetMetric, err)
			context.Status(http.StatusNotFound)
			return
		}
		context.String(http.StatusOK, value)
		logger.Infof("metric get success with type: %s, name: %s, value: %s", metricType, metricName, value)
	}
}
