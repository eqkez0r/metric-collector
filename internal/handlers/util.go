package handlers

import (
	"errors"
	"log"
	"strings"
)

var (
	errNotEnoughParams = errors.New("not enough params")
)

func getMetricParams(path, prefix string) ([]string, error) {
	log.Println("path", path)
	newPath := path[len(prefix):]
	log.Println("newPath", newPath)
	pathItems := strings.Split(newPath, "/")
	log.Println("pathItems", pathItems)
	if len(pathItems) != 2 {
		return nil, errNotEnoughParams
	}
	return pathItems, nil
}
