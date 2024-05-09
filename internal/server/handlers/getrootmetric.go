package handlers

import (
	"encoding/json"
	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	errPointGetRootMetrics = "error in GET /"
)

func GetRootMetricsHandler(
	logger *zap.SugaredLogger,
	storage store.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Infof("root metrics handler was called")
		var data map[string][]store.Metric
		var err error
		err = retry.Retry(logger, 3, func() error {
			data, err = storage.GetMetrics()
			return err
		})
		if err != nil {
			logger.Errorf("%s: %v", errPointGetRootMetrics, err)
			context.Status(http.StatusInternalServerError)
			return
		}
		bytes, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			logger.Errorf("%s: %v", errPointGetRootMetrics, err)
			context.Status(http.StatusInternalServerError)
			return
		}
		logger.Infof("get metrics was success")
		context.Data(http.StatusOK, "text/html", bytes)
	}
}
