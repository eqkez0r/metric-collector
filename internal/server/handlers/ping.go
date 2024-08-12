package handlers

import (
	"net/http"

	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Ping(
	logger *zap.SugaredLogger,
	conn *pgxpool.Pool,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("ping...")
		if conn == nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if err := retry.Retry(logger, 3, func() error { return conn.Ping(c) }); err != nil {
			logger.Info("ping failed...")
			c.Status(http.StatusInternalServerError)
			return
		}
		logger.Info("ping success")
		c.Status(http.StatusOK)
	}
}
