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
