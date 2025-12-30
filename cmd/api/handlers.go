package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"eventix-api/internal/models"
	"eventix-api/pkg/database"
	"eventix-api/pkg/jwt"
	"eventix-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================================================
// REQUEST/RESPONSE DTOs
// ============================================================================

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type CreateEventRequest struct {
	Title        string          `json:"title" validate:"required"`
	Description  string          `json:"description" validate:"required"`
	Category     string          `json:"category" validate:"required"`
	Location     string          `json:"location" validate:"required"`
	StartTime    time.Time       `json:"start_time" validate:"required"`
	EndTime      time.Time       `json:"end_time" validate:"required"`
	MaxAttendees int             `json:"max_attendees" validate:"required"`
	TicketTiers  []TicketTierReq `json:"ticket_tiers" validate:"required,min=1"`
}

type TicketTierReq struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
}

type ReserveTicketRequest struct {
	TierID   string `json:"tier_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1,max=10"`
}

type CreateOrderRequest struct {
	ReservationID string `json:"reservation_id" validate:"required"`
}

type ValidateQRRequest struct {
	QRCode  string `json:"qr_code" validate:"required"`
	EventID string `json:"event_id" validate:"required"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone,omitempty"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type EventResponse struct {
	ID           uuid.UUID            `json:"id"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Category     models.EventCategory `json:"category"`
	Location     string               `json:"location"`
	StartTime    time.Time            `json:"start_time"`
	EndTime      time.Time            `json:"end_time"`
	Status       models.EventStatus   `json:"status"`
	MaxAttendees int                  `json:"max_attendees"`
	OrganizerID  uuid.UUID            `json:"organizer_id"`
	TicketsSold  int                  `json:"tickets_sold"`
	TicketTiers  []TicketTierResponse `json:"ticket_tiers,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
}

type TicketTierResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Sold        int       `json:"sold"`
	Available   int       `json:"available"`
}

type TicketResponse struct {
	ID         uuid.UUID           `json:"id"`
	EventID    uuid.UUID           `json:"event_id"`
	EventTitle string              `json:"event_title"`
	TierName   string              `json:"tier_name"`
	QRCode     string              `json:"qr_code"`
	Status     models.TicketStatus `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
}

type OrderResponse struct {
	ID          uuid.UUID          `json:"id"`
	TotalAmount float64            `json:"total_amount"`
	Status      models.OrderStatus `json:"status"`
	TicketCount int                `json:"ticket_count"`
	CreatedAt   time.Time          `json:"created_at"`
}

// ============================================================================
// AUTH HANDLERS
// ============================================================================

// RegisterHandler godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} object{success=bool,message=string,data=UserResponse}
// @Failure 400 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 409 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /auth/register [post]
func RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if !utils.IsValidEmail(req.Email) {
		return utils.BadRequestResponse(c, "Invalid email format")
	}

	if !utils.IsValidPassword(req.Password) {
		return utils.BadRequestResponse(c, "Password must be at least 8 characters with uppercase, lowercase, and number")
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var existingUser models.User
	result := database.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return utils.ConflictResponse(c, "Email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to process registration")
	}

	user := models.User{
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Phone:         req.Phone,
		Role:          models.RoleAttendee,
		EmailVerified: false,
		IsActive:      true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create user")
	}

	userResponse := UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		Role:          string(user.Role),
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User registered successfully",
		"data":    userResponse,
	})
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} object{success=bool,data=TokenResponse}
// @Failure 400 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /auth/login [post]
func LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return utils.UnauthorizedResponse(c, "Invalid email or password")
	}

	if !user.IsActive {
		return utils.UnauthorizedResponse(c, "Account is deactivated")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return utils.UnauthorizedResponse(c, "Invalid email or password")
	}

	tokenPair, err := jwt.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate tokens")
	}

	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	response := TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresAt - time.Now().Unix(),
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// RefreshTokenHandler godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} object{success=bool,data=object{access_token=string,token_type=string,expires_in=int}}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /auth/refresh [post]
func RefreshTokenHandler(c *fiber.Ctx) error {
	var req RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	claims, err := jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid or expired refresh token")
	}

	var user models.User
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid token")
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		return utils.UnauthorizedResponse(c, "User not found")
	}

	if !user.IsActive {
		return utils.UnauthorizedResponse(c, "Account is deactivated")
	}

	tokenPair, err := jwt.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate tokens")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"token_type":    "Bearer",
			"expires_in":    tokenPair.ExpiresAt - time.Now().Unix(),
		},
	})
}

// ============================================================================
// USER HANDLERS
// ============================================================================

// GetCurrentUserHandler godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} object{success=bool,data=UserResponse}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /users/me [get]
func GetCurrentUserHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	id, err := uuid.Parse(userID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user ID")
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	userResponse := UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		Role:          string(user.Role),
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    userResponse,
	})
}

// ============================================================================
// EVENT HANDLERS
// ============================================================================

// ListEventsHandler godoc
// @Summary List all events
// @Description Get a list of all published events
// @Tags Events
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query string false "Filter by category"
// @Param status query string false "Filter by status"
// @Success 200 {object} object{success=bool,data=[]EventResponse,pagination=object{page=int,limit=int,total=int}}
// @Router /events [get]
func ListEventsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	category := c.Query("category")
	status := c.Query("status", string(models.EventPublished))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := database.DB.Model(&models.Event{}).Preload("TicketTiers")

	// Apply filters
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var events []models.Event
	if err := query.Offset(offset).Limit(limit).Order("start_time ASC").Find(&events).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch events")
	}

	eventResponses := make([]EventResponse, len(events))
	for i, event := range events {
		tierResponses := make([]TicketTierResponse, len(event.TicketTiers))
		for j, tier := range event.TicketTiers {
			tierResponses[j] = TicketTierResponse{
				ID:          tier.ID,
				Name:        tier.TierName,
				Description: tier.Description,
				Price:       tier.Price,
				Quantity:    tier.TotalQuantity,
				Sold:        tier.TotalQuantity - tier.AvailableQuantity,
				Available:   tier.AvailableQuantity,
			}
		}

		eventResponses[i] = EventResponse{
			ID:           event.ID,
			Title:        event.Title,
			Description:  event.Description,
			Category:     event.Category,
			Location:     event.Location,
			StartTime:    event.StartTime,
			EndTime:      event.EndTime,
			Status:       event.Status,
			MaxAttendees: 0,
			OrganizerID:  event.OrganizerID,
			// TicketsSold not in model
			TicketTiers: tierResponses,
			CreatedAt:   event.CreatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    eventResponses,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetEventHandler godoc
// @Summary Get event by ID
// @Description Get detailed information about a specific event
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} object{success=bool,data=EventResponse}
// @Failure 404 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /events/{id} [get]
func GetEventHandler(c *fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid event ID")
	}

	var event models.Event
	if err := database.DB.Preload("TicketTiers").First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NotFoundResponse(c, "Event not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to fetch event")
	}

	tierResponses := make([]TicketTierResponse, len(event.TicketTiers))
	for i, tier := range event.TicketTiers {
		tierResponses[i] = TicketTierResponse{
			ID:          tier.ID,
			Name:        tier.TierName,
			Description: tier.Description,
			Price:       tier.Price,
			Quantity:    tier.TotalQuantity,
			Sold:        tier.TotalQuantity - tier.AvailableQuantity,
			Available:   tier.AvailableQuantity,
		}
	}

	eventResponse := EventResponse{
		ID:           event.ID,
		Title:        event.Title,
		Description:  event.Description,
		Category:     event.Category,
		Location:     event.Location,
		StartTime:    event.StartTime,
		EndTime:      event.EndTime,
		Status:       event.Status,
		MaxAttendees: 0,
		OrganizerID:  event.OrganizerID,
		// TicketsSold not in model
		TicketTiers: tierResponses,
		CreatedAt:   event.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    eventResponse,
	})
}

// CreateEventHandler godoc
// @Summary Create a new event
// @Description Create a new event (Organizer/Admin only)
// @Tags Events
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param event body CreateEventRequest true "Event details"
// @Success 201 {object} object{success=bool,message=string,data=EventResponse}
// @Failure 400 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 403 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /events [post]
func CreateEventHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	// Check if user has organizer role
	if role != string(models.RoleOrganizer) && role != string(models.RoleAdmin) {
		return utils.ForbiddenResponse(c, "Only organizers and admins can create events")
	}

	var req CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validation
	if req.EndTime.Before(req.StartTime) {
		return utils.BadRequestResponse(c, "End time must be after start time")
	}
	if req.StartTime.Before(time.Now()) {
		return utils.BadRequestResponse(c, "Event start time cannot be in the past")
	}
	if len(req.TicketTiers) == 0 {
		return utils.BadRequestResponse(c, "At least one ticket tier is required")
	}

	// Get user's organizer profile
	var organizer models.Organizer
	uid, _ := uuid.Parse(userID)
	if err := database.DB.Where("user_id = ?", uid).First(&organizer).Error; err != nil {
		// If no organizer profile, create one
		organizer = models.Organizer{
			UserID:             uid,
			OrganizationName:   "Personal", // Default
			VerificationStatus: models.VerificationPending,
		}
		if err := database.DB.Create(&organizer).Error; err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to create organizer profile")
		}
	}

	// Create event
	event := models.Event{
		Title:       req.Title,
		Description: req.Description,
		Category:    models.EventCategory(req.Category),
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      models.EventDraft,
		// MaxAttendees not in model
		OrganizerID: organizer.ID,
		// TicketsSold not in model
	}

	// Start transaction
	tx := database.DB.Begin()
	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		return utils.InternalServerErrorResponse(c, "Failed to create event")
	}

	// Create ticket tiers
	for _, tierReq := range req.TicketTiers {
		tier := models.TicketTier{
			EventID:       event.ID,
			TierName:      tierReq.Name,
			Description:   tierReq.Description,
			Price:         tierReq.Price,
			TotalQuantity: tierReq.Quantity, AvailableQuantity: tierReq.Quantity,
			// Sold calculated from TotalQuantity - AvailableQuantity
		}
		if err := tx.Create(&tier).Error; err != nil {
			tx.Rollback()
			return utils.InternalServerErrorResponse(c, "Failed to create ticket tiers")
		}
	}

	tx.Commit()

	// Reload event with tiers
	database.DB.Preload("TicketTiers").First(&event, event.ID)

	tierResponses := make([]TicketTierResponse, len(event.TicketTiers))
	for i, tier := range event.TicketTiers {
		tierResponses[i] = TicketTierResponse{
			ID:          tier.ID,
			Name:        tier.TierName,
			Description: tier.Description,
			Price:       tier.Price,
			Quantity:    tier.TotalQuantity,
			Sold:        tier.TotalQuantity - tier.AvailableQuantity,
			Available:   tier.TotalQuantity - tier.TotalQuantity - tier.AvailableQuantity,
		}
	}

	eventResponse := EventResponse{
		ID:           event.ID,
		Title:        event.Title,
		Description:  event.Description,
		Category:     event.Category,
		Location:     event.Location,
		StartTime:    event.StartTime,
		EndTime:      event.EndTime,
		Status:       event.Status,
		MaxAttendees: 0,
		OrganizerID:  event.OrganizerID,
		// TicketsSold not in model
		TicketTiers: tierResponses,
		CreatedAt:   event.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Event created successfully",
		"data":    eventResponse,
	})
}

// ============================================================================
// TICKET HANDLERS
// ============================================================================

// ReserveTicketHandler godoc
// @Summary Reserve a ticket
// @Description Reserve a ticket for an event (15-minute hold)
// @Tags Tickets
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param request body ReserveTicketRequest true "Ticket reservation details"
// @Success 200 {object} object{success=bool,message=string,data=object{reservation_id=string,expires_at=string}}
// @Failure 400 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /tickets/reserve [post]
func ReserveTicketHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req ReserveTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	tierID, err := uuid.Parse(req.TierID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tier ID")
	}

	// Get ticket tier
	var tier models.TicketTier
	if err := database.DB.Preload("Event").First(&tier, tierID).Error; err != nil {
		return utils.NotFoundResponse(c, "Ticket tier not found")
	}

	// Check availability
	available := tier.TotalQuantity - tier.TotalQuantity - tier.AvailableQuantity
	if available < req.Quantity {
		return utils.BadRequestResponse(c, fmt.Sprintf("Only %d tickets available", available))
	}

	// Check event status
	if tier.Event.Status != models.EventPublished {
		return utils.BadRequestResponse(c, "Event is not available for booking")
	}

	// Create reservation (simplified - in production, use Redis TTL)
	reservationID := uuid.New()
	expiresAt := time.Now().Add(15 * time.Minute)

	// Store reservation in cache (placeholder - implement with Redis)
	// For now, return reservation details
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tickets reserved successfully",
		"data": fiber.Map{
			"reservation_id": reservationID.String(),
			"tier_id":        tier.ID.String(),
			"quantity":       req.Quantity,
			"unit_price":     tier.Price,
			"total_price":    tier.Price * float64(req.Quantity),
			"expires_at":     expiresAt.Format(time.RFC3339),
			"user_id":        userID,
		},
	})
}

// GetMyTicketsHandler godoc
// @Summary Get user's tickets
// @Description Get all tickets owned by the authenticated user
// @Tags Tickets
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} object{success=bool,data=[]TicketResponse}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /tickets/my-tickets [get]
func GetMyTicketsHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	uid, _ := uuid.Parse(userID)
	var tickets []models.Ticket
	if err := database.DB.Where("owner_id = ?", uid).
		Preload("Tier").
		Preload("Tier.Event").
		Order("created_at DESC").
		Find(&tickets).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch tickets")
	}

	ticketResponses := make([]TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketResponses[i] = TicketResponse{
			ID:         ticket.ID,
			EventID:    ticket.Tier.EventID,
			EventTitle: ticket.Tier.Event.Title,
			TierName:   ticket.Tier.TierName,
			QRCode:     ticket.QRCode,
			Status:     ticket.Status,
			CreatedAt:  ticket.CreatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    ticketResponses,
	})
}

// ============================================================================
// ORDER HANDLERS
// ============================================================================

// CreateOrderHandler godoc
// @Summary Create an order
// @Description Create an order for reserved tickets
// @Tags Orders
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param order body CreateOrderRequest true "Order details"
// @Success 201 {object} object{success=bool,message=string,data=OrderResponse}
// @Failure 400 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /orders [post]
func CreateOrderHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// In a real app, validate reservation from cache
	// For now, create a simplified order
	uid, _ := uuid.Parse(userID)

	order := models.Order{
		UserID:      uid,
		TotalAmount: 0, // Will be calculated based on tickets
		Status:      models.OrderPending,
	}

	// This is simplified - in production, create from reservation
	if err := database.DB.Create(&order).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create order")
	}

	orderResponse := OrderResponse{
		ID:          order.ID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		TicketCount: 0,
		CreatedAt:   order.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Order created successfully",
		"data":    orderResponse,
	})
}

// GetMyOrdersHandler godoc
// @Summary Get user's orders
// @Description Get all orders placed by the authenticated user
// @Tags Orders
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} object{success=bool,data=[]OrderResponse}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /orders/my-orders [get]
func GetMyOrdersHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	uid, _ := uuid.Parse(userID)
	var orders []models.Order
	if err := database.DB.Where("user_id = ?", uid).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch orders")
	}

	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = OrderResponse{
			ID:          order.ID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			TicketCount: 0, // Count tickets
			CreatedAt:   order.CreatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    orderResponses,
	})
}

// ============================================================================
// CHECKIN HANDLERS
// ============================================================================

// ValidateQRCodeHandler godoc
// @Summary Validate QR code
// @Description Validate a ticket QR code for event check-in (Organizer/Admin only)
// @Tags Check-in
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param request body ValidateQRRequest true "QR validation details"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 403 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /checkin/validate [post]
func ValidateQRCodeHandler(c *fiber.Ctx) error {
	var req ValidateQRRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid event ID")
	}

	// Find ticket by QR code
	var ticket models.Ticket
	if err := database.DB.Where("qr_code = ?", req.QRCode).
		Preload("Tier").
		Preload("Tier.Event").
		First(&ticket).Error; err != nil {
		return utils.NotFoundResponse(c, "Invalid QR code")
	}

	// Verify event matches
	if ticket.Tier.EventID != eventID {
		return utils.BadRequestResponse(c, "QR code is not for this event")
	}

	// Check if already checked in
	if ticket.Status == models.TicketUsed {
		return utils.BadRequestResponse(c, "Ticket already checked in")
	}

	// Check if ticket is valid
	if ticket.Status != models.TicketActive {
		return utils.BadRequestResponse(c, fmt.Sprintf("Ticket is %s", ticket.Status))
	}

	// Update ticket status
	ticket.Status = models.TicketUsed
	if err := database.DB.Save(&ticket).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update ticket")
	}

	// Create check-in record
	checkin := models.Checkin{
		TicketID:  ticket.ID,
		ScannedAt: time.Now(),
		ScannedBy: uuid.MustParse(c.Locals("user_id").(string)),
	}
	database.DB.Create(&checkin)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Check-in successful",
		"data": fiber.Map{
			"ticket_id":  ticket.ID,
			"status":     ticket.Status,
			"scanned_at": checkin.ScannedAt,
		},
	})
}

// ============================================================================
// ADMIN HANDLERS
// ============================================================================

// GetAdminStatsHandler godoc
// @Summary Get admin statistics
// @Description Get platform statistics and metrics (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} object{success=bool,data=object}
// @Failure 401 {object} object{success=bool,error=object{code=string,message=string}}
// @Failure 403 {object} object{success=bool,error=object{code=string,message=string}}
// @Router /admin/stats [get]
func GetAdminStatsHandler(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != string(models.RoleAdmin) {
		return utils.ForbiddenResponse(c, "Admin access required")
	}

	var totalUsers int64
	var totalEvents int64
	var totalTickets int64
	var totalRevenue float64

	database.DB.Model(&models.User{}).Count(&totalUsers)
	database.DB.Model(&models.Event{}).Count(&totalEvents)
	database.DB.Model(&models.Ticket{}).Count(&totalTickets)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderPaid).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_users":        totalUsers,
			"total_events":       totalEvents,
			"total_tickets_sold": totalTickets,
			"total_revenue":      totalRevenue,
		},
	})
}
