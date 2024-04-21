package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Logger(
	logger *zap.SugaredLogger,
) gin.HandlerFunc {
	logFn := func(context *gin.Context) {

		start := time.Now()

		respData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: context.Writer,
			resData:        respData,
		}

		context.Writer = &lw

		context.Next()

		respData.status = context.Writer.Status()
		respData.size = lw.Size()

		duration := time.Since(start)

		logger.Infoln(
			"URL", context.Request.URL,
			"METHOD", context.Request.Method,
			"STATUS", respData.status,
			"SIZE", respData.size,
			"DURATION", duration,
		)
	}
	return gin.HandlerFunc(logFn)
}

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
