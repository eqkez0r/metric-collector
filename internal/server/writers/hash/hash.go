package hashwriter

import (
	"encoding/base64"
	"github.com/Eqke/metric-collector/utils/hash"
	"github.com/gin-gonic/gin"
)

type HashResponseWriter struct {
	gin.ResponseWriter
	Key string
}

func New(key string) *HashResponseWriter {
	return &HashResponseWriter{Key: key}
}

func (w *HashResponseWriter) Write(data []byte) (int, error) {
	if w.Key != "" {
		h := hash.Hash(data, w.Key)
		sign := base64.StdEncoding.EncodeToString(h)
		w.ResponseWriter.Header().Add("HashSHA256", sign)
	}

	return w.ResponseWriter.Write(data)
}
