package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
)

func SubnetTrust(
	logger *zap.SugaredLogger,
	trustedSubnet string,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		if trustedSubnet != "" {
			headerIP := c.GetHeader("X-Real-IP")
			logger.Infof("received request from IP %s", headerIP)
			_, ipnet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				logger.Warnf("failed to parse trusted subnet %s: %s", trustedSubnet, err)
			}
			ip := net.ParseIP(headerIP)
			if !ipnet.Contains(ip) {
				logger.Errorf("ip: %s is not trusted", headerIP)
				c.AbortWithStatus(http.StatusForbidden)
			}
		}

		c.Next()
	}
}
