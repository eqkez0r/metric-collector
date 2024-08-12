package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	return func(c *gin.Context) {
		logger.Info("/updates: recieving metrics batch")
		if c.ContentType() != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPostMetricUpdates, c.ContentType())
			c.Status(http.StatusBadRequest)
			return
		}
		arr := []metric.Metrics{}
		if err := c.BindJSON(&arr); err != nil {
			logger.Errorf("%s: %v", errPointPostMetricJSON, err)
			c.Status(http.StatusBadRequest)
			return
		}

		logger.Info("metrics batch was recieved")
		if err := retry.Retry(logger, 3, func() error {
			return s.SetMetrics(c, arr)
		}); err != nil {
			logger.Errorf("%s: %v", err, storage.ErrIsUnknownType)
			if errors.Is(err, storage.ErrIsUnknownType) {
				c.Status(http.StatusBadRequest)
				return
			}
			if errors.Is(err, storage.ErrIDIsEmpty) {
				c.Status(http.StatusNotFound)
				return
			}
			if errors.Is(err, storage.ErrValueIsEmpty) {
				c.Status(http.StatusNotFound)
				return
			}
			c.Status(http.StatusInternalServerError)
			return
		}
		logger.Info("metrics batch was saved")
		c.Status(http.StatusOK)
	}
}
