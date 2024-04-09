package handlers

import (
	"testing"
)

func TestGetRootMetricsHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
