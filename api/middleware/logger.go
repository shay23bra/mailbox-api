package middleware

import (
	"time"

	"mailbox-api/logger"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				l.Error("Request error",
					"method", c.Request.Method,
					"path", path,
					"query", query,
					"status", c.Writer.Status(),
					"latency", latency,
					"error", e,
				)
			}
		} else {
			l.Info("Request processed",
				"method", c.Request.Method,
				"path", path,
				"query", query,
				"status", c.Writer.Status(),
				"latency", latency,
			)
		}
	}
}
