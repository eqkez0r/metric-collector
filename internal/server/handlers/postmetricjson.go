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
	errPointPostMetricJSON = "error in POST /update: "
)

//go:generate moq -out newJSONMetricProvider_moq_test.go . NewJSONMetricProvider
type NewJSONMetricProvider interface {
	SetMetric(context.Context, metric.Metrics) error
}

func POSTMetricJSONHandler(
	logger *zap.SugaredLogger,
	s NewJSONMetricProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Infof("/update: recieving metrics batch")
		if c.ContentType() != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPointPostMetricJSON, c.ContentType())
			c.Status(http.StatusBadRequest)
			return
		}
		var m metric.Metrics
		if err := c.BindJSON(&m); err != nil {
			logger.Errorf("%s: %v", errPointPostMetricJSON, err)
			c.Status(http.StatusBadRequest)
			return
		}

		logger.Info("metric was received", m)

		if err := retry.Retry(logger, 3, func() error {
			return s.SetMetric(c, m)
		}); err != nil {
			logger.Errorf("%s: %v", errPointPostMetricJSON, err)
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
		logger.Infof("metric was saved with type: %s, name: %s",
			m.MType, m.ID)
		c.Status(http.StatusOK)
	}
}
