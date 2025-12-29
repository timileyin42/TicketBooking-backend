package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TicketStatus represents ticket status
type TicketStatus string

const (
	TicketReserved  TicketStatus = "reserved"
	TicketActive    TicketStatus = "active"
	TicketUsed      TicketStatus = "used"
	TicketCancelled TicketStatus = "cancelled"
	TicketRefunded  TicketStatus = "refunded"
)

// Ticket represents a ticket
type Ticket struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TierID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"tier_id"`
	OrderID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	QRCode      string         `gorm:"uniqueIndex;not null" json:"qr_code"`
	Status      TicketStatus   `gorm:"type:varchar(20);default:'reserved';index" json:"status"`
	CheckedInAt *time.Time     `json:"checked_in_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Tier    TicketTier `gorm:"foreignKey:TierID" json:"tier,omitempty"`
	Order   Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Owner   User       `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Checkin *Checkin   `gorm:"foreignKey:TicketID" json:"checkin,omitempty"`
}

// BeforeCreate sets the ID before creating
func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// IsValid checks if ticket is valid for use
func (t *Ticket) IsValid() bool {
	return t.Status == TicketActive && t.CheckedInAt == nil
}

// OrderStatus represents order status
type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderFailed    OrderStatus = "failed"
	OrderCancelled OrderStatus = "cancelled"
	OrderRefunded  OrderStatus = "refunded"
)

// Order represents an order
type Order struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalAmount float64        `gorm:"not null" json:"total_amount"`
	Currency    string         `gorm:"default:'USD'" json:"currency"`
	Status      OrderStatus    `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tickets  []Ticket  `gorm:"foreignKey:OrderID" json:"tickets,omitempty"`
	Payments []Payment `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

// BeforeCreate sets the ID before creating
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentAuthorized PaymentStatus = "authorized"
	PaymentProcessing PaymentStatus = "processing"
	PaymentCompleted  PaymentStatus = "completed"
	PaymentFailed     PaymentStatus = "failed"
	PaymentRefunding  PaymentStatus = "refunding"
	PaymentRefunded   PaymentStatus = "refunded"
)

// PaymentProvider represents payment providers
type PaymentProvider string

const (
	ProviderPaystack PaymentProvider = "paystack"
	ProviderStripe   PaymentProvider = "stripe"
)

// Payment represents a payment
type Payment struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"order_id"`
	Provider        PaymentProvider `gorm:"type:varchar(20);not null" json:"provider"`
	Amount          float64         `gorm:"not null" json:"amount"`
	Currency        string          `gorm:"default:'USD'" json:"currency"`
	TransactionID   string          `gorm:"uniqueIndex" json:"transaction_id"`
	PaymentIntentID string          `json:"payment_intent_id,omitempty"`
	Status          PaymentStatus   `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	Metadata        string          `gorm:"type:jsonb" json:"metadata,omitempty"`
	PaidAt          *time.Time      `json:"paid_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Relationships
	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// BeforeCreate sets the ID before creating
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Checkin represents a ticket check-in
type Checkin struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TicketID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"ticket_id"`
	EventID    uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	ScannedBy  uuid.UUID `gorm:"type:uuid;not null" json:"scanned_by"`
	ScannedAt  time.Time `gorm:"not null;index" json:"scanned_at"`
	Location   string    `json:"location,omitempty"`
	DeviceInfo string    `json:"device_info,omitempty"`

	// Relationships
	Ticket  Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	Event   Event  `gorm:"foreignKey:EventID" json:"event,omitempty"`
	Scanner User   `gorm:"foreignKey:ScannedBy" json:"scanner,omitempty"`
}

// BeforeCreate sets the ID before creating
func (c *Checkin) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// NotificationType represents notification types
type NotificationType string

const (
	NotificationEmail NotificationType = "email"
	NotificationPush  NotificationType = "push"
	NotificationSMS   NotificationType = "sms"
)

// NotificationChannel represents notification channels
type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelPush  NotificationChannel = "push"
	ChannelSMS   NotificationChannel = "sms"
)

// Notification represents a notification
type Notification struct {
	ID        uuid.UUID           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID           `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      NotificationType    `gorm:"type:varchar(20);not null" json:"type"`
	Channel   NotificationChannel `gorm:"type:varchar(20);not null" json:"channel"`
	Subject   string              `json:"subject"`
	Message   string              `gorm:"type:text;not null" json:"message"`
	Metadata  string              `gorm:"type:jsonb" json:"metadata,omitempty"`
	IsRead    bool                `gorm:"default:false" json:"is_read"`
	SentAt    *time.Time          `json:"sent_at,omitempty"`
	CreatedAt time.Time           `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate sets the ID before creating
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
