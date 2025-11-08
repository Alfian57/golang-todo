package middleware

import (
	"time"

	"github.com/Alfian57/golang-todo/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ZapLogger creates a gin middleware for logging HTTP requests
// It accepts a logger.Logger interface for dependency injection
func ZapLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Skip logging for health check and swagger endpoints (optional)
		if path == "/health" || path == "/ping" {
			return
		}

		// Build log fields
		fields := []logger.Field{
			logger.F("method", method),
			logger.F("path", path),
			logger.F("query", query),
			logger.F("status", statusCode),
			logger.F("latency", latency),
			logger.F("ip", clientIP),
			logger.F("user_agent", c.Request.UserAgent()),
		}

		// Development: Log all requests
		// Production: Only log error (4xx, 5xx) or slow requests
		if gin.Mode() == gin.ReleaseMode {
			// Production mode
			if statusCode >= 400 {
				// Error requests
				if errorMessage != "" {
					fields = append(fields, logger.F("error", errorMessage))
				}
				if statusCode >= 500 {
					log.Error("HTTP Request Error", fields...)
				} else {
					log.Warn("HTTP Request Client Error", fields...)
				}
			} else if latency > 1*time.Second {
				// Slow requests (> 1 second)
				log.Warn("HTTP Slow Request", fields...)
			}
			// Success requests (2xx, 3xx) not logged in production
		} else {
			// Development mode - log all requests
			if statusCode >= 500 {
				if errorMessage != "" {
					fields = append(fields, logger.F("error", errorMessage))
				}
				log.Error("HTTP Request", fields...)
			} else if statusCode >= 400 {
				log.Warn("HTTP Request", fields...)
			} else {
				log.Info("HTTP Request", fields...)
			}
		}
	}
}

// ZapRecovery creates a gin middleware for recovering from panics
// It accepts a logger.Logger interface for dependency injection
func ZapRecovery(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered",
					logger.F("error", err),
					logger.F("path", c.Request.URL.Path),
					logger.F("method", c.Request.Method),
					logger.F("ip", c.ClientIP()),
				)

				// Return 500 error
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
