package middleware

import (
	"kbt/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		utils.HTTPRequestsInFlight.Inc()
		c.Next()
		utils.HTTPRequestsInFlight.Dec()

		path := c.FullPath()
		if path == "" {
			path = "/"
		}

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		utils.HttpRequestsTotal.WithLabelValues(method, path, status).Inc()

		utils.HttpRequestDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}
}
