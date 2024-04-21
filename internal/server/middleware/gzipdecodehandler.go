package middleware

import (
	"compress/gzip"
	"github.com/Eqke/metric-collector/internal/compress"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	errPointGzipDecode = "error in gzip decode handler: "
)

var avaliableTypes = map[string]bool{
	"text/html":        true,
	"application/json": true,
	"html/text":        true,
}

func GzipMiddleware(
	logger *zap.SugaredLogger,
) gin.HandlerFunc {
	return func(context *gin.Context) {
		ae := context.GetHeader("Accept-Encoding")
		ct := context.GetHeader("Content-Type")
		ac := context.GetHeader("Accept")

		if (strings.Contains(ae, "gzip")) &&
			(avaliableTypes[ct] && avaliableTypes[ac]) {
			logger.Debug("request was encoded. decoding...")
			logger.Debugf("Content-Type: %s, Accept: %s, Accept-Encoding: %s", ct, ac, ae)
			w := context.Writer
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()
			w = compress.GzipWriter{
				ResponseWriter: context.Writer,
				Writer:         gzipWriter,
			}
			context.Header("Content-Encoding", "application/gzip")
		}

		ce := context.GetHeader("Content-Encoding")
		logger.Debugf("Content-Encoding: %s", ce)
		if strings.Contains(ce, "gzip") {
			gzipReader, err := gzip.NewReader(context.Request.Body)
			if err != nil {
				logger.Errorf("%s: %v", errPointGzipDecode, err)
				context.Status(http.StatusInternalServerError)
				return
			}
			defer gzipReader.Close()

			context.Request.Body = gzipReader
		}
		context.Next()
	}
}
