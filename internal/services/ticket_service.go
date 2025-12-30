package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"eventix-api/internal/models"
	"eventix-api/pkg/cache"
	"eventix-api/pkg/database"
	"eventix-api/pkg/utils"
)

// ReservationData represents a ticket reservation in Redis
type ReservationData struct {
	ReservationID string    `json:"reservation_id"`
	UserID        uuid.UUID `json:"user_id"`
	TierID        uuid.UUID `json:"tier_id"`
	EventID       uuid.UUID `json:"event_id"`
	Quantity      int       `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	TotalPrice    float64   `json:"total_price"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// TicketService handles ticket-related operations
type TicketService struct{}

// NewTicketService creates a new ticket service
func NewTicketService() *TicketService {
	return &TicketService{}
}

// CreateReservation creates a temporary ticket reservation
func (s *TicketService) CreateReservation(userID, tierID uuid.UUID, quantity int) (*ReservationData, error) {
	// Get ticket tier
	var tier models.TicketTier
	if err := database.DB.Preload("Event").First(&tier, tierID).Error; err != nil {
		return nil, fmt.Errorf("tier not found: %w", err)
	}

	// Check availability
	if tier.AvailableQuantity < quantity {
		return nil, fmt.Errorf("only %d tickets available", tier.AvailableQuantity)
	}

	// Check event status
	if tier.Event.Status != models.EventPublished {
		return nil, fmt.Errorf("event is not available for booking")
	}

	// Create reservation
	reservation := &ReservationData{
		ReservationID: utils.GenerateReservationID(),
		UserID:        userID,
		TierID:        tierID,
		EventID:       tier.EventID,
		Quantity:      quantity,
		UnitPrice:     tier.Price,
		TotalPrice:    utils.CalculateTotalPrice(tier.Price, quantity),
		ExpiresAt:     time.Now().Add(utils.ReservationExpirySeconds()),
		CreatedAt:     time.Now(),
	}

	// Store in Redis
	key := utils.GetReservationKey(reservation.ReservationID)
	data, _ := json.Marshal(reservation)

	ctx := context.Background()
	if err := cache.Client.Set(ctx, key, data, utils.ReservationExpirySeconds()).Err(); err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// Temporarily reduce available quantity
	tier.AvailableQuantity -= quantity
	if err := database.DB.Save(&tier).Error; err != nil {
		// Rollback Redis
		cache.Client.Del(ctx, key)
		return nil, fmt.Errorf("failed to reserve tickets: %w", err)
	}

	return reservation, nil
}

// GetReservation retrieves a reservation from Redis
func (s *TicketService) GetReservation(reservationID string) (*ReservationData, error) {
	key := utils.GetReservationKey(reservationID)
	ctx := context.Background()

	data, err := cache.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("reservation not found or expired")
	}

	var reservation ReservationData
	if err := json.Unmarshal([]byte(data), &reservation); err != nil {
		return nil, fmt.Errorf("invalid reservation data")
	}

	return &reservation, nil
}

// DeleteReservation deletes a reservation and releases tickets
func (s *TicketService) DeleteReservation(reservationID string) error {
	reservation, err := s.GetReservation(reservationID)
	if err != nil {
		return err // Already expired or not found
	}

	// Release tickets back to available quantity
	if err := database.DB.Model(&models.TicketTier{}).
		Where("id = ?", reservation.TierID).
		UpdateColumn("available_quantity", gorm.Expr("available_quantity + ?", reservation.Quantity)).
		Error; err != nil {
		return fmt.Errorf("failed to release tickets: %w", err)
	}

	// Delete from Redis
	ctx := context.Background()
	cache.Client.Del(ctx, utils.GetReservationKey(reservationID))

	return nil
}

// CreateTicketsFromOrder creates tickets for a paid order
func (s *TicketService) CreateTicketsFromOrder(orderID, tierID, userID uuid.UUID, quantity int) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, quantity)

	for i := 0; i < quantity; i++ {
		ticketID := uuid.New()
		qrCode := utils.GenerateQRCode(ticketID)

		tickets[i] = models.Ticket{
			ID:      ticketID,
			TierID:  tierID,
			OrderID: orderID,
			OwnerID: userID,
			QRCode:  qrCode,
			Status:  models.TicketActive,
		}
	}

	// Create all tickets in database
	if err := database.DB.Create(&tickets).Error; err != nil {
		return nil, fmt.Errorf("failed to create tickets: %w", err)
	}

	return tickets, nil
}

// ProcessOrderPayment processes order and creates tickets
func (s *TicketService) ProcessOrderPayment(orderID uuid.UUID) error {
	// Get order with tickets
	var order models.Order
	if err := database.DB.Preload("Tickets").First(&order, orderID).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Update order status to paid
	order.Status = models.OrderPaid
	if err := database.DB.Save(&order).Error; err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Create payment record
	payment := models.Payment{
		OrderID:  orderID,
		Amount:   order.TotalAmount,
		Currency: order.Currency,
		Provider: models.ProviderPaystack, // Default
		Status:   models.PaymentCompleted,
	}
	now := time.Now()
	payment.PaidAt = &now

	if err := database.DB.Create(&payment).Error; err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// ValidateTicketForCheckin validates a ticket for check-in
func (s *TicketService) ValidateTicketForCheckin(qrCode string, eventID uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := database.DB.Where("qr_code = ?", qrCode).
		Preload("Tier").
		Preload("Tier.Event").
		First(&ticket).Error; err != nil {
		return nil, fmt.Errorf("invalid QR code")
	}

	// Verify event matches
	if ticket.Tier.EventID != eventID {
		return nil, fmt.Errorf("QR code is not for this event")
	}

	// Check if already checked in
	if ticket.Status == models.TicketUsed {
		return nil, fmt.Errorf("ticket already checked in")
	}

	// Check if ticket is active
	if ticket.Status != models.TicketActive {
		return nil, fmt.Errorf("ticket is %s", ticket.Status)
	}

	return &ticket, nil
}

// CheckInTicket marks a ticket as checked in
func (s *TicketService) CheckInTicket(ticket *models.Ticket, validatorID uuid.UUID, eventID uuid.UUID) (*models.Checkin, error) {
	// Update ticket status
	ticket.Status = models.TicketUsed
	now := time.Now()
	ticket.CheckedInAt = &now

	if err := database.DB.Save(ticket).Error; err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	// Create check-in record
	checkin := models.Checkin{
		TicketID:  ticket.ID,
		EventID:   eventID,
		ScannedBy: validatorID,
		ScannedAt: now,
	}

	if err := database.DB.Create(&checkin).Error; err != nil {
		return nil, fmt.Errorf("failed to create check-in record: %w", err)
	}

	return &checkin, nil
}
