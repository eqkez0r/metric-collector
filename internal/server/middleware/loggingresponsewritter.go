package middleware

import "github.com/gin-gonic/gin"

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		gin.ResponseWriter
		resData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.resData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.resData.status = statusCode // захватываем код статуса
}
