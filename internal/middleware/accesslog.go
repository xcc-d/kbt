package middleware

import (
	"time"

	"kbt/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		user := ""

		if claims := GetClaims(c); claims != nil {
			user = claims.Username
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Int64("latency_ms", latency.Milliseconds()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user", user),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		switch {
		case statusCode >= 500:
			utils.L().Error("access", fields...)
		case statusCode >= 400:
			utils.L().Warn("access", fields...)
		default:
			utils.L().Info("access", fields...)
		}
	}
}
