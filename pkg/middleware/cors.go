package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"eventix-api/pkg/config"
	"eventix-api/pkg/logger"
	"go.uber.org/zap"
)

// CORS creates a CORS middleware
func CORS(cfg *config.CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.AllowedOrigins, ","),
		AllowMethods:     joinStrings(cfg.AllowedMethods, ","),
		AllowHeaders:     joinStrings(cfg.AllowedHeaders, ","),
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Type",
		MaxAge:           86400,
	})
}

// Recover creates a panic recovery middleware
func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			logger.Error("Panic recovered",
				zap.Any("error", e),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			)
		},
	})
}

func joinStrings(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += sep + slice[i]
	}
	return result
}
