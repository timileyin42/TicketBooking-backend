package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateQRCode generates a unique QR code for a ticket
func GenerateQRCode(ticketID uuid.UUID) string {
	// Format: TICKET-{UUID}-{TIMESTAMP}-{RANDOM}
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomStr := base64.URLEncoding.EncodeToString(randomBytes)[:8]

	return fmt.Sprintf("TICKET-%s-%d-%s", ticketID.String()[:8], timestamp, randomStr)
}

// ValidateQRCodeFormat validates QR code format
func ValidateQRCodeFormat(qrCode string) bool {
	// Simple validation - check if it starts with TICKET-
	return len(qrCode) > 15 && qrCode[:7] == "TICKET-"
}

// GenerateReservationID generates a unique reservation ID
func GenerateReservationID() string {
	return uuid.New().String()
}

// CalculateTotalPrice calculates total price for tickets
func CalculateTotalPrice(price float64, quantity int) float64 {
	return price * float64(quantity)
}

// GetReservationKey returns Redis key for reservation
func GetReservationKey(reservationID string) string {
	return fmt.Sprintf("reservation:%s", reservationID)
}

// ReservationExpirySeconds returns reservation expiry duration (15 minutes)
func ReservationExpirySeconds() time.Duration {
	return 15 * time.Minute
}
