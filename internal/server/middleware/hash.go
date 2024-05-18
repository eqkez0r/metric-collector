package middleware

import (
	"bytes"
	"encoding/base64"
	hashwriter "github.com/Eqke/metric-collector/internal/server/writers/hash"
	"github.com/Eqke/metric-collector/utils/hash"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	errPointHash = "error in hash handler: "
)

func Hash(
	logger *zap.SugaredLogger,
	hashKey string,
) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Infof("hash header: %s", context.GetHeader("HashSHA256"))
		if hashKey != "" && context.GetHeader("HashSHA256") != "" {
			header := context.GetHeader("HashSHA256")
			logger.Info("checking hash")
			body, err := io.ReadAll(context.Request.Body)
			if err != nil {
				logger.Errorf("%s: %v", errPointHash, err)
				context.Status(http.StatusInternalServerError)
				return
			}
			context.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			//logger.Infof("data: %s", body)
			//logger.Infof("hashKey %s", hashKey)
			h := hash.Hash(body, hashKey)
			sign := base64.StdEncoding.EncodeToString(h)
			if header != sign {
				logger.Errorf("%s: %s != %s", errPointHash, header, sign)
				context.Status(http.StatusBadRequest)
				return
			}
			//logger.Infof("hash: %s hash calculated: %s", header, sign)
			logger.Info("hash checked successful")
			context.Writer = &hashwriter.HashResponseWriter{
				ResponseWriter: context.Writer,
				Key:            hashKey,
			}
		}

		context.Next()

	}
}

func HashMiddleware(
	logger *zap.SugaredLogger,
	hashKey string,
) gin.HandlerFunc {
	return Hash(logger, hashKey)
}
