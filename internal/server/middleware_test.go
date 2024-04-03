package httpserver

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type handler struct {
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func Test_middleware(t *testing.T) {
	type want struct {
		statusCode  int
		method      string
		contentType string
	}

	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			name:   "Success request",
			url:    "http://localhost:8080/update/gauge/test/12.32",
			method: http.MethodPost,
			want: want{
				statusCode:  http.StatusOK,
				method:      http.MethodPost,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Invalid request",
			url:    "http://localhost:8080/update/k/test/12.32",
			method: http.MethodGet,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				method:      http.MethodPost,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			middleware(handler{}).ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
