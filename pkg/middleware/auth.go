package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"eventix-api/pkg/jwt"
	"eventix-api/pkg/utils"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "Authorization header required")
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.UnauthorizedResponse(c, "Invalid authorization header format")
		}

		token := parts[1]

		// Validate token
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			return utils.UnauthorizedResponse(c, "Invalid or expired token")
		}

		// Set user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return utils.UnauthorizedResponse(c, "User not authenticated")
		}

		userRole := role.(string)

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return utils.ForbiddenResponse(c, "You don't have permission to access this resource")
	}
}

// OptionalAuthMiddleware validates JWT token if present
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			return c.Next()
		}

		// Set user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) string {
	userID := c.Locals("user_id")
	if userID == nil {
		return ""
	}
	return userID.(string)
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *fiber.Ctx) string {
	email := c.Locals("email")
	if email == nil {
		return ""
	}
	return email.(string)
}

// GetUserRole extracts user role from context
func GetUserRole(c *fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(c *fiber.Ctx) bool {
	return c.Locals("user_id") != nil
}
