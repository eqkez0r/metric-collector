package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"log"
	"net/http"
	"strconv"
)

const GaugePath = "/update/gauge/"

type GaugeHandler struct {
	Storage *localstorage.LocalStorage
}

func (gh GaugeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("gauge handler was called")

	path := r.URL.Path

	pathItems, err := getMetricParams(path, GaugePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	name := pathItems[0]
	value := pathItems[1]
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = gh.Storage.Gauge(name, metric.Gauge(valueFloat)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
