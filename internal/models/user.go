package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents user roles in the system
type UserRole string

const (
	RoleAttendee  UserRole = "attendee"
	RoleOrganizer UserRole = "organizer"
	RoleAdmin     UserRole = "admin"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	FirstName     string         `gorm:"not null" json:"first_name"`
	LastName      string         `gorm:"not null" json:"last_name"`
	Phone         string         `json:"phone,omitempty"`
	Role          UserRole       `gorm:"type:varchar(20);not null;default:'attendee'" json:"role"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// OAuth fields
	OAuthProvider string `json:"oauth_provider,omitempty"`
	OAuthID       string `json:"oauth_id,omitempty"`

	// Relationships
	Organizer *Organizer `gorm:"foreignKey:UserID" json:"organizer,omitempty"`
	Orders    []Order    `gorm:"foreignKey:UserID" json:"-"`
	Tickets   []Ticket   `gorm:"foreignKey:OwnerID" json:"-"`
}

// BeforeCreate sets the ID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// VerificationStatus represents organizer verification status
type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationApproved VerificationStatus = "approved"
	VerificationRejected VerificationStatus = "rejected"
)

// Organizer represents an event organizer
type Organizer struct {
	ID                 uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	OrganizationName   string             `gorm:"not null" json:"organization_name"`
	Description        string             `gorm:"type:text" json:"description"`
	Website            string             `json:"website,omitempty"`
	Logo               string             `json:"logo,omitempty"`
	VerificationStatus VerificationStatus `gorm:"type:varchar(20);default:'pending'" json:"verification_status"`
	VerifiedAt         *time.Time         `json:"verified_at,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`

	// Relationships
	User   User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Events []Event `gorm:"foreignKey:OrganizerID" json:"-"`
}

// BeforeCreate sets the ID before creating
func (o *Organizer) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// EventStatus represents event status
type EventStatus string

const (
	EventDraft       EventStatus = "draft"
	EventUnderReview EventStatus = "under_review"
	EventPublished   EventStatus = "published"
	EventActive      EventStatus = "active"
	EventCompleted   EventStatus = "completed"
	EventCancelled   EventStatus = "cancelled"
)

// EventCategory represents event categories
type EventCategory string

const (
	CategoryMusic      EventCategory = "music"
	CategorySports     EventCategory = "sports"
	CategoryArts       EventCategory = "arts"
	CategoryTechnology EventCategory = "technology"
	CategoryBusiness   EventCategory = "business"
	CategoryEducation  EventCategory = "education"
	CategoryOther      EventCategory = "other"
)

// Event represents an event
type Event struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizerID uuid.UUID      `gorm:"type:uuid;not null;index" json:"organizer_id"`
	Title       string         `gorm:"not null" json:"title"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Category    EventCategory  `gorm:"type:varchar(50);not null" json:"category"`
	Location    string         `gorm:"not null" json:"location"`
	Venue       string         `json:"venue"`
	StartTime   time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime     time.Time      `gorm:"not null" json:"end_time"`
	BannerURL   string         `json:"banner_url"`
	Status      EventStatus    `gorm:"type:varchar(20);default:'draft';index" json:"status"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Organizer   Organizer    `gorm:"foreignKey:OrganizerID" json:"organizer,omitempty"`
	TicketTiers []TicketTier `gorm:"foreignKey:EventID" json:"ticket_tiers,omitempty"`
	Checkins    []Checkin    `gorm:"foreignKey:EventID" json:"-"`
}

// BeforeCreate sets the ID before creating
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// IsActive checks if event is currently active
func (e *Event) IsActive() bool {
	now := time.Now()
	return e.Status == EventActive && now.After(e.StartTime) && now.Before(e.EndTime)
}

// TicketTier represents a ticket tier for an event
type TicketTier struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EventID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"event_id"`
	TierName          string         `gorm:"not null" json:"tier_name"`
	Description       string         `gorm:"type:text" json:"description"`
	Price             float64        `gorm:"not null" json:"price"`
	Currency          string         `gorm:"default:'USD'" json:"currency"`
	TotalQuantity     int            `gorm:"not null" json:"total_quantity"`
	AvailableQuantity int            `gorm:"not null" json:"available_quantity"`
	SaleStartTime     *time.Time     `json:"sale_start_time,omitempty"`
	SaleEndTime       *time.Time     `json:"sale_end_time,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Event   Event    `gorm:"foreignKey:EventID" json:"event,omitempty"`
	Tickets []Ticket `gorm:"foreignKey:TierID" json:"-"`
}

// BeforeCreate sets the ID before creating
func (t *TicketTier) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// IsAvailable checks if tickets are available for sale
func (t *TicketTier) IsAvailable() bool {
	now := time.Now()

	if t.AvailableQuantity <= 0 {
		return false
	}

	if t.SaleStartTime != nil && now.Before(*t.SaleStartTime) {
		return false
	}

	if t.SaleEndTime != nil && now.After(*t.SaleEndTime) {
		return false
	}

	return true
}

// TableName specifies the table name for TicketTier
func (TicketTier) TableName() string {
	return "ticket_tiers"
}
