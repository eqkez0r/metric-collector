package handlers

import (
	"context"
	"net/http"

	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PingProvider interface {
	Ping(context.Context) error
}

func Ping(
	logger *zap.SugaredLogger,
	store PingProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("ping...")
		if err := retry.Retry(logger, 3, func() error {
			return store.Ping(c)
		}); err != nil {
			logger.Info("ping failed...")
			c.Status(http.StatusInternalServerError)
			return
		}
		logger.Info("ping success")
		c.Status(http.StatusOK)
	}
}
