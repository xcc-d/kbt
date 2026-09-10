package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(allowOrigins []string) gin.HandlerFunc {
	allAllow := len(allowOrigins) == 1 && allowOrigins[0] == "*"
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allAllow {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" && contains(allowOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "7200")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
