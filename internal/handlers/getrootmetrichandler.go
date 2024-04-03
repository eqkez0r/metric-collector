package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetRootMetricsHandler(storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		const errPoint = "error in GET /"
		metrics, err := storage.GetMetrics()
		if err != nil {
			log.Println(errPoint, err)
			context.Status(http.StatusInternalServerError)
			return
		}
		log.Println("metrics got")
		context.IndentedJSON(http.StatusOK, metrics)
	}
}
