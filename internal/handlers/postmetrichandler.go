package handlers

import (
	"errors"
	"github.com/Eqke/metric-collector/internal/storage"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	errPointPostMetric = "error in POST /update/:type/:name/:value"
)

func POSTMetricHandler(storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {

		metricType := context.Param("type")
		if metricType != metric.TypeCounter.String() && metricType != metric.TypeGauge.String() {
			log.Println(e.WrapError(errPointPostMetric, errors.New("invalid metric type "+metricType)))
			context.Status(http.StatusBadRequest)
			return
		}
		metricName := context.Param("name")
		if metricName == "" {
			log.Println(e.WrapError(errPointPostMetric, errors.New("empty metric name")))
			context.Status(http.StatusNotFound)
			return
		}
		metricValue := context.Param("value")
		log.Println("metric was received", metricType, metricName, metricValue)
		err := storage.SetValue(metricType, metricName, metricValue)
		if err != nil {
			log.Println(e.WrapError(errPointPostMetric, err))
			context.Status(http.StatusBadRequest)
			return
		}
		log.Println("metric was updated")
		context.Status(http.StatusOK)

	}
}
