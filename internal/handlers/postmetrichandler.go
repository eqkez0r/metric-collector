package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func POSTMetricHandler(storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		const errPoint = "error in POST /update/:type/:name/:value"
		metricType := context.Param("type")
		if metricType != metric.TypeCounter.String() && metricType != metric.TypeGauge.String() {
			log.Println(errPoint, "unsupported metric type", metricType)
			context.Status(http.StatusBadRequest)
			return
		}
		metricName := context.Param("name")
		if metricName == "" {
			log.Println(errPoint, "empty metric name")
			context.Status(http.StatusNotFound)
			return
		}
		metricValue := context.Param("value")
		log.Println("metric was received", metricType, metricName, metricValue)
		err := storage.SetValue(metricType, metricName, metricValue)
		if err != nil {
			log.Println(errPoint, err)
			context.Status(http.StatusBadRequest)
			return
		}
		log.Println("metric was updated")
		context.Status(http.StatusOK)

	}
}
