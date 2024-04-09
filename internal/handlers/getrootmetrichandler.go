package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	errPointGetRootMetrics = "error in GET /"
)

func GetRootMetricsHandler(storage storage.Storage) gin.HandlerFunc {
	return func(context *gin.Context) {

		metrics, err := storage.GetMetrics()
		if err != nil {
			log.Println(e.WrapError(errPointGetRootMetrics, err))
			context.Status(http.StatusInternalServerError)
			return
		}
		log.Println("metrics got")
		context.IndentedJSON(http.StatusOK, metrics)
	}
}
