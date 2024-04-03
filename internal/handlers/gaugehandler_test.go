package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGaugeHandler_ServeHTTP(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Success request",
			url:  "http://localhost:8080/update/gauge/test/12.32",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Without metric name",
			url:  "http://localhost:8080/update/gauge/12.32",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Without value",
			url:  "http://localhost:8080/update/gauge/test/",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "with invalid value",
			url:  "http://localhost:8080/update/gauge/test/abc",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()

			var gauge GaugeHandler
			gauge.Storage = localstorage.New()
			gauge.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)

		})
	}
}
