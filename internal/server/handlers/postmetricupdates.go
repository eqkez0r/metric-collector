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
	errPostMetricUpdates = "error in POST /updates: " // error in POST /update
)

//go:generate moq -out batchMetricProvider_moq_test.go . BatchMetricProvider
type BatchMetricProvider interface {
	SetMetrics(context.Context, []metric.Metrics) error
}

func PostMetricUpdates(
	logger *zap.SugaredLogger,
	s BatchMetricProvider) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Infof("/updates: recieving metrics batch")
		if context.ContentType() != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPostMetricUpdates, context.ContentType())
			context.Status(http.StatusBadRequest)
			return
		}
		arr := []metric.Metrics{}
		if err := context.BindJSON(&arr); err != nil {
			logger.Errorf("%s: %v", errPointPostMetricJSON, err)
			context.Status(http.StatusBadRequest)
			return
		}

		logger.Infof("metrics batch was recieved")
		if err := retry.Retry(logger, 3, func() error {
			return s.SetMetrics(context, arr)
		}); err != nil {
			logger.Errorf("%s: %v", err, storage.ErrIsUnknownType)
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
		logger.Infof("metrics batch was saved")
		context.Status(http.StatusOK)
	}
}
