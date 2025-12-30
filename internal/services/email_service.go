package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"

	"eventix-api/pkg/cache"
	"eventix-api/pkg/config"
	"eventix-api/pkg/utils"
)

// EmailService handles email operations
type EmailService struct {
	client    *resend.Client
	fromEmail string
	fromName  string
	cfg       *config.EmailConfig
	templates map[string]*template.Template
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.EmailConfig) *EmailService {
	client := resend.NewClient(cfg.ResendAPIKey)
	service := &EmailService{
		client:    client,
		fromEmail: cfg.FromEmail,
		fromName:  cfg.FromName,
		cfg:       cfg,
		templates: make(map[string]*template.Template),
	}

	// Load email templates
	service.loadTemplates()

	return service
}

// loadTemplates loads all email templates from the templates/email directory
func (s *EmailService) loadTemplates() {
	templateDir := "templates/email"

	templates := map[string]string{
		"verify_email":       "verify_email.html",
		"welcome":            "welcome.html",
		"order_confirmation": "order_confirmation.html",
	}

	for name, filename := range templates {
		templatePath := filepath.Join(templateDir, filename)

		// Check if file exists
		if _, err := os.Stat(templatePath); err != nil {
			// Template file doesn't exist, skip
			continue
		}

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			// Failed to parse template, skip
			continue
		}

		s.templates[name] = tmpl
	}
}

// renderTemplate renders an email template with the given data
func (s *EmailService) renderTemplate(templateName string, data interface{}) (string, error) {
	tmpl, exists := s.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(userID uuid.UUID, email, firstName, frontendURL string) error {
	// Generate verification token
	token := utils.GenerateReservationID()

	// Store in Redis with 24 hour expiry
	key := fmt.Sprintf("email_verify:%s", token)
	ctx := context.Background()
	if err := cache.Client.Set(ctx, key, userID.String(), 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	// Create verification link
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)

	// Prepare template data
	data := map[string]interface{}{
		"FirstName":        firstName,
		"VerificationLink": verificationLink,
	}

	// Render template
	htmlBody, err := s.renderTemplate("verify_email", data)
	if err != nil {
		return err
	}

	// Send email using Resend
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{email},
		Subject: "Verify Your Email - Eventix",
		Html:    htmlBody,
	}

	_, err = s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// VerifyEmailToken verifies an email verification token
func (s *EmailService) VerifyEmailToken(token string) (uuid.UUID, error) {
	key := fmt.Sprintf("email_verify:%s", token)
	ctx := context.Background()

	userIDStr, err := cache.Client.Get(ctx, key).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid or expired verification token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token data")
	}

	// Delete token after use (one-time use)
	cache.Client.Del(ctx, key)

	return userID, nil
}

// SendOrderConfirmationEmail sends order confirmation with tickets
func (s *EmailService) SendOrderConfirmationEmail(email, firstName string, orderID uuid.UUID, totalAmount float64, ticketCount int) error {
	// Prepare template data
	data := map[string]interface{}{
		"FirstName":   firstName,
		"OrderID":     orderID.String(),
		"TicketCount": ticketCount,
		"TotalAmount": fmt.Sprintf("%.2f", totalAmount),
	}

	// Render template
	htmlBody, err := s.renderTemplate("order_confirmation", data)
	if err != nil {
		return err
	}

	// Send email
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{email},
		Subject: "Order Confirmed - Your Tickets are Ready!",
		Html:    htmlBody,
	}

	_, err = s.client.Emails.Send(params)
	return err
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(email, firstName string) error {
	// Prepare template data
	data := map[string]interface{}{
		"FirstName": firstName,
	}

	// Render template
	htmlBody, err := s.renderTemplate("welcome", data)
	if err != nil {
		return err
	}

	// Send email
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{email},
		Subject: "Welcome to Eventix! 🎉",
		Html:    htmlBody,
	}

	_, err = s.client.Emails.Send(params)
	return err
}
