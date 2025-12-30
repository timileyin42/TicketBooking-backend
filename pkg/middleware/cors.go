package middleware

import (
	"eventix-api/pkg/config"
	"eventix-api/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

// CORS creates a CORS middleware
func CORS(cfg *config.CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*", // Allow all origins in development
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowCredentials: false, // Must be false when using wildcard origin
		ExposeHeaders:    "Content-Length,Content-Type,Authorization",
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
