package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func GetMetricJSONHandler(
	logger *zap.SugaredLogger,
	s storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		var m metric.Metrics
		if err := context.BindJSON(&m); err != nil {
			logger.Errorf("%s: %v", errPointGetMetric, err)
			context.Status(http.StatusBadRequest)
			return
		}
		finVal, err := s.GetMetric(m)
		if err != nil {
			logger.Errorf("%s: %v", errPointGetMetric, err)
			context.Status(http.StatusBadRequest)
			return
		}

		context.Header("Content-Type", "application/json")
		context.JSON(http.StatusOK, finVal)
	}
}
