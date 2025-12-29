package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"eventix-api/pkg/cache"
	"eventix-api/pkg/config"
	"eventix-api/pkg/database"
	"eventix-api/pkg/jwt"
	"eventix-api/pkg/logger"
	"eventix-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
	"go.uber.org/zap"

	// Import domain handlers
)

// @title Eventix Ticket Booking API
// @version 1.0
// @description Production-grade ticket booking system API for events, organizers, and attendees
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@eventix.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.App.LogLevel, cfg.App.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Eventix API Server",
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Environment),
	)

	// Initialize JWT
	jwt.Init(&cfg.JWT)

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Connect to Redis
	if err := cache.Connect(&cfg.Redis); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer cache.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
		BodyLimit:    10 * 1024 * 1024, // 10MB
	})

	// Global middleware
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS(&cfg.CORS))
	app.Use(middleware.RateLimiter(&cfg.Limits))

	// Swagger documentation endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/health", healthCheckHandler)

	// API routes
	api := app.Group(fmt.Sprintf("/api/%s", cfg.App.Version))

	setupRoutes(api, cfg)

	// 404 handler
	app.Use(notFoundHandler)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Server starting", zap.String("address", addr))

	if err := app.Listen(addr); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

func setupRoutes(api fiber.Router, cfg *config.Config) {
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Register endpoint - to be implemented"})
	})
	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Login endpoint - to be implemented"})
	})
	auth.Post("/refresh", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Refresh token endpoint - to be implemented"})
	})

	// Protected routes
	protected := api.Group("", middleware.AuthMiddleware())

	// User routes
	users := protected.Group("/users")
	users.Get("/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get current user - to be implemented"})
	})

	// Event routes
	events := api.Group("/events")
	events.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "List events - to be implemented"})
	})
	events.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get event - to be implemented"})
	})

	organizerEvents := protected.Group("/events", middleware.RoleMiddleware("organizer", "admin"))
	organizerEvents.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Create event - to be implemented"})
	})

	// Ticket routes
	tickets := protected.Group("/tickets")
	tickets.Post("/reserve", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Reserve ticket - to be implemented"})
	})
	tickets.Get("/my-tickets", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get my tickets - to be implemented"})
	})

	// Order routes
	orders := protected.Group("/orders")
	orders.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Create order - to be implemented"})
	})
	orders.Get("/my-orders", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get my orders - to be implemented"})
	})

	// Check-in routes
	checkin := protected.Group("/checkin", middleware.RoleMiddleware("organizer", "admin"))
	checkin.Post("/validate", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Validate QR code - to be implemented"})
	})

	// Admin routes
	admin := protected.Group("/admin", middleware.RoleMiddleware("admin"))
	admin.Get("/stats", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get admin stats - to be implemented"})
	})

	logger.Info("Routes registered successfully")
}

func healthCheckHandler(c *fiber.Ctx) error {
	// Check database health
	dbHealth := "healthy"
	if err := database.Health(); err != nil {
		dbHealth = "unhealthy"
	}

	// Check cache health
	cacheHealth := "healthy"
	if err := cache.Health(); err != nil {
		cacheHealth = "unhealthy"
	}

	return c.JSON(fiber.Map{
		"status":   "ok",
		"service":  "eventix-api",
		"database": dbHealth,
		"cache":    cacheHealth,
	})
}

func notFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "NOT_FOUND",
			"message": "The requested resource was not found",
		},
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logger.Error("Request error",
		zap.Int("status_code", code),
		zap.String("path", c.Path()),
		zap.Error(err),
	)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": err.Error(),
		},
	})
}
