package middleware

import (
	"bytes"
	"github.com/Eqke/metric-collector/internal/server/httpserver/writers"
	"io"
	"net/http"

	h "github.com/Eqke/metric-collector/utils/hash"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	errPointHash = "error in hash middleware: "
)

func Hash(
	logger *zap.SugaredLogger,
	hashKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hashHeader := c.GetHeader("HashSHA256")
		logger.Infof("%s: %s", "hash header", hashHeader)

		if hashHeader != "" && hashKey != "" {
			logger.Info("Checking hash")
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Errorf("%s: %v", errPointHash, err)
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			sign := h.Sign(body, hashKey)
			if sign != hashHeader {
				logger.Errorf("%s: %v", errPointHash, "hash not equal")
				c.Status(http.StatusBadRequest)
				return
			}
			logger.Infof("%s: %s", "hash checked successfully", hashHeader)
			w := writers.NewHashWriter(c.Writer, logger, hashKey)
			c.Writer = w
		}

		c.Next()
	}
}
