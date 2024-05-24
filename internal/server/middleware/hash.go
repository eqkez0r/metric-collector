package middleware

import (
	"bytes"
	"github.com/Eqke/metric-collector/internal/server/writers"
	hash2 "github.com/Eqke/metric-collector/utils/hash"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	errPointHash = "error in hash middleware: "
)

func Hash(
	logger *zap.SugaredLogger,
	hashKey string) gin.HandlerFunc {
	return func(context *gin.Context) {
		hashHeader := context.GetHeader("HashSHA256")
		logger.Infof("%s: %s", "hash header", hashHeader)
		//if hashHeader == "" {
		//	context.Next()
		//	return
		//}
		if hashHeader != "" && hashKey != "" {
			logger.Infof("Checking hash")
			body, err := io.ReadAll(context.Request.Body)
			if err != nil {
				logger.Errorf("%s: %v", errPointHash, err)
				context.Status(http.StatusInternalServerError)
				return
			}

			context.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			sign := hash2.Sign(body, hashKey)
			if sign != hashHeader {
				logger.Errorf("%s: %v", errPointHash, "hash not equal")
				context.Status(http.StatusBadRequest)
				return
			}
			logger.Infof("%s: %s", "hash checked successfully", hashHeader)
			w := writers.NewHashWriter(context.Writer, logger, hashKey)
			context.Writer = w
		}

		context.Next()
	}
}
