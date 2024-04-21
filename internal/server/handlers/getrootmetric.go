package handlers

import (
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
		context.IndentedJSON(http.StatusOK, metrics)
	}
}
