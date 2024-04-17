package handlers

import (
	ls "github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounterHandler_ServeHTTP(t *testing.T) {
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
			url:  "http://localhost:8080/update/counter/test/10",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Without metric name",
			url:  "http://localhost:8080/update/counter/10",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Without value",
			url:  "http://localhost:8080/update/counter/test",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "with invalid value",
			url:  "http://localhost:8080/update/counter/test/abc",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := zaptest.NewLogger(t).Sugar()
			localstorage := ls.New(l)
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			gin.DisableConsoleColor()
			gin.SetMode(gin.ReleaseMode)

			engine := gin.New()
			engine.POST("/update/:type/:name/:value", POSTMetricHandler(l, localstorage))
			engine.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
