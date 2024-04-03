package handlers

import (
	"log"
	"net/http"
	"strings"
)

const UpdatePath = "/update/"

type UpdateHandler struct{}

var (
	SupportedMetricsType = map[string]bool{
		"counter": true,
		"gauge":   true,
	}
)

func (gh UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	path = path[len(UpdatePath):]
	log.Println(path)
	pathItems := strings.Split(path, "/")
	if v, ok := SupportedMetricsType[pathItems[0]]; !ok || !v {
		log.Println("Unsupported metric type", ok, v)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
