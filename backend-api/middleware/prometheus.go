package middleware

import (
	"strconv"
	"time"

	"livecode-api/internal/metrics"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/metrics" || path == "/health" {
			c.Next()
			return
		}

		startTime := time.Now()

		metrics.HTTPRequestsInFlight.Inc()

		c.Next()

		duration := time.Since(startTime).Seconds()

		method := c.Request.Method
		path = c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := strconv.Itoa(c.Writer.Status())

		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HTTPRequestsInFlight.Dec()
	}
}
