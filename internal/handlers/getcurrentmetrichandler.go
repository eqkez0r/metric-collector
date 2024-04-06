package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	errPointGetMetric = "error in GET /value/:type/:name"
)

func GETMetricHandler(storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {

		metricType := context.Param("type")
		metricName := context.Param("name")

		value, err := storage.GetValue(metricType, metricName)
		if err != nil {
			log.Println(e.WrapError(errPointGetMetric, err))
			context.Status(http.StatusNotFound)
			return
		}
		context.String(http.StatusOK, value)
	}
}
