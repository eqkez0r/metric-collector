package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	errPointGetMetricJSON = "error in POST /value: "
)

func GetMetricJSONHandler(
	logger *zap.SugaredLogger,
	s storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {

		var m metric.Metrics
		ct := context.GetHeader("Content-Type")
		context.Header("Content-Type", "application/json")
		if ct != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPointGetMetricJSON, ct)
			context.Status(http.StatusNotFound)
			return
		}
		if err := context.BindJSON(&m); err != nil {
			logger.Errorf("%s: %v", errPointGetMetricJSON, err)
			context.Status(http.StatusBadRequest)
			return
		}
		if m.ID == "" {
			logger.Errorf("%s: empty metric name", errPointGetMetricJSON)
			context.Status(http.StatusNotFound)
			return
		}
		finVal, err := s.GetMetric(m)
		if err != nil {
			logger.Errorf("%s: %v", errPointGetMetricJSON, err)
			context.Status(http.StatusNotFound)
			return
		}
		logger.Infof("metric get success with type: %s, name: %s", m.MType, m.ID)
		context.JSON(http.StatusOK, finVal)
	}
}
