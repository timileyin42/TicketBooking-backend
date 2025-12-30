package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Payment  PaymentConfig
	Kafka    KafkaConfig
	RabbitMQ RabbitMQConfig
	S3       S3Config
	Email    EmailConfig
	Server   ServerConfig
	Limits   LimitsConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
	LogLevel    string
	LogFormat   string
}

type DatabaseConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConnections int
	MaxIdleConns   int
	MaxLifetime    time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	TTL      time.Duration
}

type JWTConfig struct {
	Secret             string
	Expiry             time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

type PaymentConfig struct {
	PaystackSecretKey string
	PaystackPublicKey string
	StripeSecretKey   string
	StripePublicKey   string
	WebhookSecret     string
}

type KafkaConfig struct {
	Brokers            []string
	GroupID            string
	TopicPayments      string
	TopicNotifications string
	TopicCheckins      string
}

type RabbitMQConfig struct {
	URL                string
	Exchange           string
	QueuePayments      string
	QueueNotifications string
}

type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	UseSSL          bool
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	SendGridKey  string
	ResendAPIKey string
}

type ServerConfig struct {
	Port              int
	FrontendURL       string
	AdminURL          string
	PrometheusPort    int
	PrometheusEnabled bool
}

type LimitsConfig struct {
	RateLimitRequests        int
	RateLimitWindow          time.Duration
	TicketReservationTimeout time.Duration
	MaxTicketsPerOrder       int
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Eventix"),
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("API_VERSION", "v1"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			LogFormat:   getEnv("LOG_FORMAT", "json"),
		},
		Database: DatabaseConfig{
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnvAsInt("DB_PORT", 5432),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", "postgres"),
			Name:           getEnv("DB_NAME", "ticket_booking"),
			SSLMode:        getEnv("DB_SSL_MODE", "disable"),
			MaxConnections: getEnvAsInt("DB_MAX_CONNECTIONS", 100),
			MaxIdleConns:   getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 10),
			MaxLifetime:    getEnvAsDuration("DB_MAX_LIFETIME", 3600*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			TTL:      getEnvAsDuration("REDIS_TTL", 3600*time.Second),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-secret-key"),
			Expiry:             getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
			RefreshTokenExpiry: getEnvAsDuration("REFRESH_TOKEN_EXPIRY", 168*time.Hour),
			Issuer:             getEnv("JWT_ISSUER", "eventix-api"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		},
		Payment: PaymentConfig{
			PaystackSecretKey: getEnv("PAYSTACK_SECRET_KEY", ""),
			PaystackPublicKey: getEnv("PAYSTACK_PUBLIC_KEY", ""),
			StripeSecretKey:   getEnv("STRIPE_SECRET_KEY", ""),
			StripePublicKey:   getEnv("STRIPE_PUBLIC_KEY", ""),
			WebhookSecret:     getEnv("WEBHOOK_SECRET", ""),
		},
		Kafka: KafkaConfig{
			Brokers:            getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			GroupID:            getEnv("KAFKA_GROUP_ID", "ticket-booking-consumer"),
			TopicPayments:      getEnv("KAFKA_TOPIC_PAYMENTS", "payments"),
			TopicNotifications: getEnv("KAFKA_TOPIC_NOTIFICATIONS", "notifications"),
			TopicCheckins:      getEnv("KAFKA_TOPIC_CHECKINS", "checkins"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:                getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Exchange:           getEnv("RABBITMQ_EXCHANGE", "ticket_booking"),
			QueuePayments:      getEnv("RABBITMQ_QUEUE_PAYMENTS", "payments_queue"),
			QueueNotifications: getEnv("RABBITMQ_QUEUE_NOTIFICATIONS", "notifications_queue"),
		},
		S3: S3Config{
			Bucket:          getEnv("S3_BUCKET", "eventix-uploads"),
			Region:          getEnv("S3_REGION", "us-east-1"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			Endpoint:        getEnv("S3_ENDPOINT", ""),
			UseSSL:          getEnvAsBool("S3_USE_SSL", true),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("SMTP_FROM_EMAIL", "noreply@eventix.com"),
			FromName:     getEnv("SMTP_FROM_NAME", "Eventix"),
			SendGridKey:  getEnv("SENDGRID_API_KEY", ""),
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		},
		Server: ServerConfig{
			Port:              getEnvAsInt("PORT", 8080),
			FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3000"),
			AdminURL:          getEnv("ADMIN_URL", "http://localhost:3001"),
			PrometheusPort:    getEnvAsInt("PROMETHEUS_PORT", 9090),
			PrometheusEnabled: getEnvAsBool("PROMETHEUS_ENABLED", true),
		},
		Limits: LimitsConfig{
			RateLimitRequests:        getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:          getEnvAsDuration("RATE_LIMIT_WINDOW", 60*time.Second),
			TicketReservationTimeout: getEnvAsDuration("TICKET_RESERVATION_TIMEOUT", 15*time.Minute),
			MaxTicketsPerOrder:       getEnvAsInt("MAX_TICKETS_PER_ORDER", 10),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key" {
		return fmt.Errorf("JWT secret must be set and cannot be default value")
	}
	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var result []string
	for _, v := range splitString(valueStr, ",") {
		if trimmed := trimSpace(v); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var result []string
	current := ""
	for _, char := range s {
		if string(char) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}

	return s[start:end]
}
