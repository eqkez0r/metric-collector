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
		data, err := storage.GetMetrics()
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
