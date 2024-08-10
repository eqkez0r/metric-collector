package handlers

import (
	"context"
	"errors"
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	errPointPostMetricJSON = "error in POST /update: "
)

//go:generate moq -out newJSONMetricProvider_moq_test.go . NewJSONMetricProvider
type NewJSONMetricProvider interface {
	SetMetric(context.Context, metric.Metrics) error
}

func POSTMetricJSONHandler(
	logger *zap.SugaredLogger,
	s NewJSONMetricProvider) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Infof("/update: recieving metrics batch")
		if context.ContentType() != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPointPostMetricJSON, context.ContentType())
			context.Status(http.StatusBadRequest)
			return
		}
		var m metric.Metrics
		if err := context.BindJSON(&m); err != nil {
			logger.Errorf("%s: %v", errPointPostMetricJSON, err)
			context.Status(http.StatusBadRequest)
			return
		}

		logger.Info("metric was received", m)

		if err := retry.Retry(logger, 3, func() error {
			return s.SetMetric(context, m)
		}); err != nil {
			logger.Errorf("%s: %v", errPointPostMetricJSON, err)
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
