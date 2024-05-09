package handlers

import (
	"github.com/Eqke/metric-collector/utils/retry"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
)

func Ping(
	logger *zap.SugaredLogger,
	conn *pgxpool.Conn,
) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Info("ping...")
		if conn == nil {
			context.Status(http.StatusInternalServerError)
			return
		}
		if err := retry.Retry(logger, 3, func() error { return conn.Ping(context) }); err != nil {
			logger.Info("ping failed...")
			context.Status(http.StatusInternalServerError)
			return
		}
		logger.Info("ping success")
		context.Status(http.StatusOK)
	}
}
