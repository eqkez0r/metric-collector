package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/Eqke/metric-collector/pkg/metrics"
	"log"
	"net/http"
	"strconv"
)

const CounterPath = "/update/counter/"

type CounterHandler struct {
	Storage *localstorage.LocalStorage
}

func (ch CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("counter handler was called")

	path := r.URL.Path

	pathItems, err := getMetricParams(path, CounterPath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	name := pathItems[0]
	value := pathItems[1]
	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = ch.Storage.Counter(name, metrics.Counter(valueInt)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

}
