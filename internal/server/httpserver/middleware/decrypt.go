package middleware

import (
	"bytes"
	"crypto/rsa"
	"github.com/Eqke/metric-collector/internal/encrypting"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func Decrypt(
	logger *zap.SugaredLogger,
	cryptoKey *rsa.PrivateKey,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Errorw("Error reading body", "error", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		decryptedBody, err := encrypting.Decrypt(cryptoKey, body)
		if err != nil {
			logger.Errorw("Error decrypting body", "error", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
		c.Next()
	}
}
