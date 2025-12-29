package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"eventix-api/pkg/logger"
	"go.uber.org/zap"
)

// Logger middleware logs HTTP requests
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Generate request ID if not present
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)
		status := c.Response().StatusCode()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		}

		// Add user ID if authenticated
		if userID := c.Locals("user_id"); userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		// Add error if present
		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		// Log based on status code
		if status >= 500 {
			logger.Error("HTTP request error", fields...)
		} else if status >= 400 {
			logger.Warn("HTTP request warning", fields...)
		} else {
			logger.Info("HTTP request", fields...)
		}

		return err
	}
}
