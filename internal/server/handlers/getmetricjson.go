package handlers

import (
	"context"
	"net/http"

	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	errPointGetMetricJSON = "error in POST /value: "
)

//go:generate moq -out jsonMetricProvider_moq_test.go . JSONMetricProvider
type JSONMetricProvider interface {
	GetMetric(context.Context, metric.Metrics) (metric.Metrics, error)
}

func GetMetricJSONHandler(
	logger *zap.SugaredLogger,
	p JSONMetricProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("/value: get metric json")
		var m metric.Metrics
		var err error
		ct := c.GetHeader("Content-Type")
		c.Header("Content-Type", "application/json")
		if ct != "application/json" {
			logger.Errorf("%s: unknown content type %s", errPointGetMetricJSON, ct)
			c.Status(http.StatusNotFound)
			return
		}
		if err = c.BindJSON(&m); err != nil {
			logger.Errorf("%s: %v", errPointGetMetricJSON, err)
			c.Status(http.StatusBadRequest)
			return
		}
		if m.ID == "" {
			logger.Errorf("%s: empty metric name", errPointGetMetricJSON)
			c.Status(http.StatusNotFound)
			return
		}
		err = retry.Retry(logger, 3, func() error {
			m, err = p.GetMetric(c, m)
			return err
		})
		if err != nil {
			logger.Errorf("%s: %v", errPointGetMetricJSON, err)
			c.Status(http.StatusNotFound)
			return
		}
		logger.Infof("metric get success with type: %s, name: %s", m.MType, m.ID)
		c.JSON(http.StatusOK, m)
	}
}
