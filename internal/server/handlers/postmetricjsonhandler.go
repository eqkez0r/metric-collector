package handlers

import (
	"errors"
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func POSTMetricJSONHandler(
	logger *zap.SugaredLogger,
	s storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.ContentType() != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPointPostMetric, context.ContentType())
			context.Status(http.StatusBadRequest)
			return
		}
		var m metric.Metrics
		if err := context.BindJSON(&m); err != nil {
			logger.Errorf("%s: %v", errPointPostMetric, err)
			context.Status(http.StatusBadRequest)
			return
		}
		logger.Info("metric was received", m)
		if err := s.SetMetric(m); err != nil {
			logger.Errorf("%s: %v", errPointPostMetric, err)
			if errors.Is(err, storage.ErrIsUnknownType) {
				context.Status(http.StatusBadRequest)
				return
			}
			if errors.Is(err, storage.ErrIDIsEmpty) {
				context.Status(http.StatusNotFound)
				return
			}
			if errors.Is(err, storage.ErrValueIsEmpty) {
				context.Status(http.StatusNotFound)
				return
			}
			context.Status(http.StatusInternalServerError)
			return
		}
		logger.Infof("metric was saved with type: %s, name: %s",
			m.MType, m.ID)
		context.Status(http.StatusOK)
	}
}
