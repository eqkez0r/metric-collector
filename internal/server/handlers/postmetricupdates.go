package handlers

import (
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

func PostMetricUpdates(
	logger *zap.SugaredLogger,
	s storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.ContentType() != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPostMetricUpdates, context.ContentType())
			context.Status(http.StatusBadRequest)
			return
		}
		arr := []metric.Metrics{}
		if err := context.BindJSON(&arr); err != nil {
			logger.Infof("getting data from body: %s", arr)
			logger.Errorf("%s: %v", errPostMetricUpdates, err)
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
