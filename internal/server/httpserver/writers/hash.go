package writers

import (
	hash2 "github.com/Eqke/metric-collector/utils/hash"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HashResponseWriter struct {
	gin.ResponseWriter
	hashKey string
	logger  *zap.SugaredLogger
}

func NewHashWriter(w gin.ResponseWriter,
	logger *zap.SugaredLogger,
	HashKey string) *HashResponseWriter {
	return &HashResponseWriter{
		ResponseWriter: w,
		hashKey:        HashKey,
		logger:         logger,
	}
}

func (w HashResponseWriter) Write(b []byte) (int, error) {
	if w.hashKey != "" {
		sign := hash2.Sign(b, w.hashKey)
		w.Header().Set("HashSHA256", sign)
	}
	return w.ResponseWriter.Write(b)
}
