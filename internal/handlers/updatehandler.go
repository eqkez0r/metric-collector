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
	newPath := path[len(UpdatePath):]
	pathItems := strings.Split(newPath, "/")
	if v, ok := SupportedMetricsType[pathItems[0]]; !ok || !v {
		log.Println("Unsupported metric type", ok, v)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
