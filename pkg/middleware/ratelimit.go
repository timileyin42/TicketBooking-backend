package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"eventix-api/pkg/config"
)

// RateLimiter creates a rate limiting middleware
func RateLimiter(cfg *config.LimitsConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.RateLimitRequests,
		Expiration: cfg.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests, please try again later",
				},
				"timestamp": time.Now().UTC(),
			})
		},
	})
}

// StrictRateLimiter creates a stricter rate limiter for sensitive endpoints
func StrictRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many attempts, please try again in a minute",
				},
				"timestamp": time.Now().UTC(),
			})
		},
	})
}
