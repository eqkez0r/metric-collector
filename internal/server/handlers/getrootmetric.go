package handlers

import (
	"encoding/json"
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	errPointGetRootMetrics = "error in GET /"
)

func GetRootMetricsHandler(
	logger *zap.SugaredLogger,
	storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Infof("root metrics handler was called")
		metrics, err := storage.GetMetrics()
		if err != nil {
			logger.Errorf("%s: %v", errPointGetRootMetrics, err)
			context.Status(http.StatusInternalServerError)
			return
		}
		logger.Infof("metrics get success")
		bytes, err := json.Marshal(metrics)
		if err != nil {
			logger.Errorf("%s: %v", errPointGetRootMetrics, err)
			context.Status(http.StatusInternalServerError)
			return
		}
		context.String(http.StatusOK, string(bytes))
		context.Header("Content-Type", "text/html")
		logger.Infof("root metrics handler was finished")
	}
}
