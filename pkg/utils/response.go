package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success   bool         `json:"success"`
	Message   string       `json:"message,omitempty"`
	Data      interface{}  `json:"data,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

// ErrorDetail represents error details
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
	Timestamp  time.Time      `json:"timestamp"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// CreatedResponse sends a created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, code, message string, details interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
	})
}

// BadRequestResponse sends a bad request error response
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "BAD_REQUEST", message, nil)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, "FORBIDDEN", message, nil)
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, "NOT_FOUND", message, nil)
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, "CONFLICT", message, nil)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *fiber.Ctx, errors interface{}) error {
	return ErrorResponse(c, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", errors)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "An internal server error occurred"
	}
	return ErrorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

// PaginatedSuccessResponse sends a paginated success response
func PaginatedSuccessResponse(c *fiber.Ctx, data interface{}, page, limit int, total int64) error {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: total,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Timestamp: time.Now().UTC(),
	})
}

// NoContentResponse sends a no content response
func NoContentResponse(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
