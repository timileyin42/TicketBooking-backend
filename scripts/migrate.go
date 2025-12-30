package main

import (
	"eventix-api/internal/models"
	"eventix-api/pkg/config"
	"eventix-api/pkg/database"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Running database migrations...")

	// Auto-migrate all models
	err = database.DB.AutoMigrate(
		&models.User{},
		&models.Organizer{},
		&models.Event{},
		&models.TicketTier{},
		&models.Ticket{},
		&models.Order{},
		&models.Payment{},
		&models.Checkin{},
		&models.Notification{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✅ All migrations completed successfully!")
}
