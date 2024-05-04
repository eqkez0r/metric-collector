package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"net/http"
)

func Ping(
	logger *zap.SugaredLogger,
	conn *pgx.Conn,
) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Info("ping...")
		if conn == nil {
			context.Status(http.StatusInternalServerError)
			return
		}
		if err := conn.Ping(context); err != nil {
			logger.Info("ping failed...")
			context.Status(http.StatusInternalServerError)
			return
		}
		logger.Info("ping success")
		context.Status(http.StatusOK)
	}
}
