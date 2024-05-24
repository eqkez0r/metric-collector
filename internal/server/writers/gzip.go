package writers

import (
	"github.com/gin-gonic/gin"
	"io"
)

type GzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
