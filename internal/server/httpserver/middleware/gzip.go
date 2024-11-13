package middleware

import (
	"compress/gzip"
	"github.com/Eqke/metric-collector/internal/server/httpserver/writers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	errPointGzipDecode = "error in gzip decode handler: "
)

var availableTypes = map[string]bool{
	"text/html":        true,
	"application/json": true,
	"html/text":        true,
}

func Gzip(
	logger *zap.SugaredLogger,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		//compress response
		ae := c.GetHeader("Accept-Encoding")
		ct := c.GetHeader("Content-Type")
		ac := c.GetHeader("Accept")
		if (strings.Contains(ae, "gzip")) &&
			(availableTypes[ct] || availableTypes[ac]) {
			w := c.Writer
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()
			wd := writers.GzipWriter{
				ResponseWriter: c.Writer,
				Writer:         gzipWriter,
			}
			c.Writer = wd
			c.Writer.Header().Set("Content-Encoding", "gzip")
		}

		//uncompressed request
		ce := c.GetHeader("Content-Encoding")
		if strings.Contains(ce, "gzip") {
			gzipReader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				logger.Errorf("%s: %v", errPointGzipDecode, err)
				c.Status(http.StatusInternalServerError)
				return
			}
			defer gzipReader.Close()

			c.Request.Body = gzipReader
		}
		c.Next()
	}
}
