package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/utils/retry"
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
	return func(c *gin.Context) {
		metricType := c.Param("type")
		metricName := c.Param("name")
		logger.Infof("metric was requested with type: %s, name: %s", metricType, metricName)
		var value string
		var err error
		err = retry.Retry(logger, 3, func() error {
			value, err = storage.GetValue(c, metricType, metricName)
			return err
		})
		if err != nil {
			logger.Errorf("%s: %v", errPointGetMetric, err)
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, value)
		logger.Infof("metric get success with type: %s, name: %s, value: %s", metricType, metricName, value)
	}
}
