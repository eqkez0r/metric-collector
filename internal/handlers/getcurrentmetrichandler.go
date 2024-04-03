package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GETMetricHandler(storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		const errPoint = "error in GET /value/:type/:name"
		metricType := context.Param("type")
		metricName := context.Param("name")

		value, err := storage.GetValue(metricType, metricName)
		if err != nil {
			log.Println(errPoint, err)
			context.Status(http.StatusNotFound)
			return
		}
		context.String(http.StatusOK, value)
	}
}
