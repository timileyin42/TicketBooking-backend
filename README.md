# 🎫 Eventix - Ticket Booking System

> A multi-platform ticket booking system enabling event discovery, secure ticket booking, and event management for attendees, organizers, and admins.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Latest-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-Latest-DC382D?style=flat&logo=redis)](https://redis.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Core Features](#core-features)
- [System Design](#system-design)
- [Database Schema](#database-schema)
- [API Design](#api-design)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Security](#security)
- [Performance](#performance)

---

## 🎯 Overview

**Eventix** is a production-grade ticket booking platform designed to solve key pain points in event discovery and ticket management:

- **Event Discovery**: Centralized platform for finding events
- **Secure Booking**: PCI-compliant payment processing with QR-code tickets
- **Real-time Validation**: QR-based check-in with offline support
- **Analytics**: Comprehensive dashboards for organizers
- **Mobile-First**: Optimized mobile experience via React Native/Flutter

### Target Users

| Role | Description |
|------|-------------|
| **Attendees** | Event-goers who browse, book, and manage tickets |
| **Organizers** | Event creators who sell tickets and track analytics |
| **Admins** | Platform operators managing users, events, and payments |

---

## 🏗️ Architecture

### System Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        WebApp[Web App<br/>React/Next.js]
        MobileApp[Mobile App<br/>React Native/Flutter]
    end
    
    subgraph "API Gateway"
        Gateway[API Gateway<br/>JWT Auth + Rate Limiting]
    end
    
    subgraph "Backend Services - Modular Monolith"
        Auth[Auth Service<br/>JWT + OAuth]
        Users[User Service]
        Events[Event Service]
        Tickets[Ticket Service<br/>QR Generation]
        Orders[Order Service]
        Payments[Payment Service<br/>Paystack/Stripe]
        Checkin[Check-in Service<br/>QR Validation]
        Notifications[Notification Service<br/>Email + Push]
        Admin[Admin Service]
    end
    
    subgraph "Data Layer"
        PostgreSQL[(PostgreSQL<br/>Primary DB)]
        Redis[(Redis<br/>Cache + Sessions)]
    end
    
    subgraph "Message Queue"
        Kafka[Kafka/RabbitMQ<br/>Event Streaming]
    end
    
    subgraph "External Services"
        PaymentGateway[Payment Gateway<br/>Paystack/Stripe]
        Storage[S3 Storage<br/>Media Assets]
        EmailProvider[Email Service<br/>SMTP/SendGrid]
    end
    
    WebApp --> Gateway
    MobileApp --> Gateway
    
    Gateway --> Auth
    Gateway --> Users
    Gateway --> Events
    Gateway --> Tickets
    Gateway --> Orders
    Gateway --> Checkin
    Gateway --> Admin
    
    Orders --> Payments
    Payments --> PaymentGateway
    
    Auth --> PostgreSQL
    Users --> PostgreSQL
    Events --> PostgreSQL
    Tickets --> PostgreSQL
    Orders --> PostgreSQL
    Checkin --> PostgreSQL
    
    Auth --> Redis
    Tickets --> Redis
    Orders --> Redis
    
    Payments --> Kafka
    Notifications --> Kafka
    Kafka --> EmailProvider
    
    Events --> Storage
    Users --> Storage
    
    style Auth fill:#4CAF50
    style Payments fill:#FF9800
    style Checkin fill:#2196F3
```

### Architecture Principles

```mermaid
mindmap
  root((Clean Architecture))
    Modular Monolith
      Domain-Driven Design
      Bounded Contexts
      Future Microservices
    Explicit Interfaces
      Repository Pattern
      Service Layer
      Handler Layer
    Testability
      Unit Tests
      Integration Tests
      Mock Interfaces
    Scalability
      Horizontal Scaling
      Event-Driven
      Caching Strategy
```

---

## 🛠️ Tech Stack

### Backend Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Language** | Go 1.22+ | High-performance backend |
| **HTTP Framework** | Fiber / net/http | RESTful API |
| **ORM** | GORM / SQLC | Database abstraction |
| **Database** | PostgreSQL | Primary data store |
| **Cache** | Redis | Session & performance |
| **Auth** | JWT + OAuth2 | Authentication & authorization |
| **Message Queue** | Kafka / RabbitMQ | Event streaming |
| **Payments** | Paystack / Stripe | Payment processing |
| **Storage** | S3-compatible | Media storage |
| **Config** | ENV / Doppler | Configuration management |
| **Migrations** | Goose | Database versioning |
| **Observability** | Prometheus + OpenTelemetry | Monitoring & tracing |

### Frontend Stack

| Platform | Technology |
|----------|------------|
| **Web** | React / Next.js |
| **Mobile** | React Native / Flutter |

---

## ✨ Core Features

### For Attendees 🎟️

```mermaid
graph LR
    A[Discover Events] --> B[Search & Filter]
    B --> C[Select Event]
    C --> D[Choose Ticket Tier]
    D --> E[Reserve Ticket]
    E --> F[Payment]
    F --> G[QR Ticket Generated]
    G --> H[Email Confirmation]
    
    style A fill:#E3F2FD
    style F fill:#FFF3E0
    style G fill:#E8F5E9
```

**Features:**
- 🔍 Event discovery & advanced search
- 💳 Secure ticket booking & payment
- 📱 QR-code based digital tickets
- 📜 Booking history & management
- 🔄 Ticket transfers
- 💸 Refund requests
- 🔔 Email & push notifications

### For Organizers 🎪

```mermaid
graph TD
    A[Create Event] --> B[Configure Details]
    B --> C[Setup Ticket Tiers]
    C --> D[Set Pricing]
    D --> E[Publish Event]
    E --> F[Monitor Sales]
    F --> G[Check-in Attendees]
    G --> H[Receive Payouts]
    
    style A fill:#FFF3E0
    style E fill:#E8F5E9
    style H fill:#E1F5FE
```

**Features:**
- 📝 Event creation & management
- 🎫 Ticket tier configuration
- 📊 Sales dashboard & analytics
- ✅ Attendee list & check-in
- 💰 Payout management
- 📈 Performance metrics

### For Admin 👨‍💼

**Features:**
- ✓ User & organizer verification
- 🔒 Event moderation & approval
- 💳 Payment oversight
- 📊 Platform metrics & reporting
- 🚫 Ban & suspension controls
- 📋 Audit logs

---

## 🔄 System Design

### Ticket Booking Flow

```mermaid
sequenceDiagram
    actor User
    participant Web as Web/Mobile App
    participant API as API Gateway
    participant Auth as Auth Service
    participant Event as Event Service
    participant Ticket as Ticket Service
    participant Order as Order Service
    participant Payment as Payment Service
    participant Queue as Message Queue
    participant Email as Email Service
    participant DB as PostgreSQL
    participant Cache as Redis
    
    User->>Web: Browse Events
    Web->>API: GET /api/v1/events
    API->>Event: Fetch Events
    Event->>Cache: Check Cache
    Cache-->>Event: Cache Hit/Miss
    alt Cache Miss
        Event->>DB: Query Events
        DB-->>Event: Return Events
        Event->>Cache: Update Cache
    end
    Event-->>Web: Event List
    
    User->>Web: Select Event & Ticket
    Web->>API: POST /api/v1/tickets/reserve
    API->>Auth: Validate Token
    Auth-->>API: Token Valid
    API->>Ticket: Reserve Ticket
    Ticket->>Cache: Set Reservation (TTL)
    Ticket-->>Web: Reservation ID
    
    User->>Web: Proceed to Payment
    Web->>API: POST /api/v1/orders
    API->>Order: Create Order
    Order->>Payment: Initialize Payment
    Payment->>DB: Save Payment Intent
    Payment-->>Web: Payment URL
    
    User->>Web: Complete Payment
    Web->>Payment: Payment Callback
    Payment->>DB: Update Payment Status
    Payment->>Ticket: Generate QR Code
    Ticket->>DB: Save Ticket
    Payment->>Queue: Publish Payment.Success
    Queue->>Email: Send Confirmation
    Email->>User: Email with QR Ticket
    Payment-->>Web: Payment Success
```

### Check-in Flow

```mermaid
sequenceDiagram
    actor Organizer
    participant Scanner as Scanner App
    participant API as API Gateway
    participant Checkin as Check-in Service
    participant Ticket as Ticket Service
    participant DB as PostgreSQL
    participant Cache as Redis
    
    Organizer->>Scanner: Scan QR Code
    Scanner->>API: POST /api/v1/checkin/validate
    API->>Checkin: Validate Ticket
    Checkin->>Cache: Check if Scanned
    
    alt Already Scanned
        Cache-->>Checkin: Duplicate Scan
        Checkin-->>Scanner: ❌ Already Used
    else Valid Ticket
        Checkin->>Ticket: Verify Ticket
        Ticket->>DB: Query Ticket
        DB-->>Ticket: Ticket Details
        
        alt Ticket Valid
            Ticket-->>Checkin: Valid
            Checkin->>DB: Mark as Checked-in
            Checkin->>Cache: Set Scanned Flag
            Checkin-->>Scanner: ✅ Success
        else Invalid Ticket
            Ticket-->>Checkin: Invalid
            Checkin-->>Scanner: ❌ Invalid Ticket
        end
    end
```

### Payment Processing

```mermaid
stateDiagram-v2
    [*] --> Pending: Order Created
    Pending --> Authorized: Payment Intent
    Authorized --> Processing: User Pays
    Processing --> Completed: Webhook Success
    Processing --> Failed: Payment Failed
    Completed --> Refunding: Refund Request
    Refunding --> Refunded: Refund Approved
    Failed --> [*]
    Refunded --> [*]
    Completed --> [*]
    
    note right of Pending
        Reserve ticket
        Start 15min timer
    end note
    
    note right of Completed
        Generate QR
        Send email
        Release to organizer
    end note
```

### Event Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Draft: Create Event
    Draft --> UnderReview: Submit for Review
    UnderReview --> Published: Admin Approves
    UnderReview --> Draft: Admin Rejects
    Published --> Active: Event Date Arrives
    Active --> Completed: Event Ends
    Published --> Cancelled: Organizer Cancels
    Completed --> [*]
    Cancelled --> [*]
    
    note right of Draft
        Editable
        Not visible
    end note
    
    note right of Published
        Public
        Ticket sales open
    end note
    
    note right of Active
        Check-in enabled
        Real-time updates
    end note
```

---

## 🗄️ Database Schema

### Core Tables

```mermaid
erDiagram
    USERS ||--o{ ORGANIZERS : "is"
    USERS ||--o{ ORDERS : "places"
    USERS ||--o{ TICKETS : "owns"
    
    ORGANIZERS ||--o{ EVENTS : "creates"
    
    EVENTS ||--o{ TICKET_TIERS : "has"
    EVENTS ||--o{ CHECKINS : "tracks"
    
    TICKET_TIERS ||--o{ TICKETS : "generates"
    
    ORDERS ||--o{ TICKETS : "contains"
    ORDERS ||--o{ PAYMENTS : "initiates"
    
    TICKETS ||--o{ CHECKINS : "validated_by"
    
    USERS {
        uuid id PK
        string email UK
        string password_hash
        string first_name
        string last_name
        string phone
        enum role
        timestamp created_at
        timestamp updated_at
    }
    
    ORGANIZERS {
        uuid id PK
        uuid user_id FK
        string organization_name
        string description
        enum verification_status
        timestamp verified_at
        timestamp created_at
    }
    
    EVENTS {
        uuid id PK
        uuid organizer_id FK
        string title
        text description
        string category
        string location
        timestamp start_time
        timestamp end_time
        string banner_url
        enum status
        timestamp created_at
        timestamp updated_at
    }
    
    TICKET_TIERS {
        uuid id PK
        uuid event_id FK
        string tier_name
        text description
        decimal price
        int total_quantity
        int available_quantity
        timestamp sale_start
        timestamp sale_end
        timestamp created_at
    }
    
    TICKETS {
        uuid id PK
        uuid tier_id FK
        uuid order_id FK
        uuid owner_id FK
        string qr_code UK
        enum status
        timestamp checked_in_at
        timestamp created_at
    }
    
    ORDERS {
        uuid id PK
        uuid user_id FK
        decimal total_amount
        enum status
        timestamp created_at
        timestamp updated_at
    }
    
    PAYMENTS {
        uuid id PK
        uuid order_id FK
        string payment_provider
        decimal amount
        string currency
        string transaction_id UK
        enum status
        timestamp paid_at
        timestamp created_at
    }
    
    CHECKINS {
        uuid id PK
        uuid ticket_id FK
        uuid event_id FK
        uuid scanned_by FK
        timestamp scanned_at
        string location
    }
    
    NOTIFICATIONS {
        uuid id PK
        uuid user_id FK
        string type
        string channel
        text message
        boolean is_read
        timestamp sent_at
        timestamp created_at
    }
```

### Database Indexes

```mermaid
graph TB
    subgraph "Performance Indexes"
        A[events.start_time<br/>events.category]
        B[tickets.qr_code<br/>tickets.owner_id]
        C[orders.user_id<br/>orders.status]
        D[payments.transaction_id<br/>payments.order_id]
    end
    
    subgraph "Unique Constraints"
        E[users.email]
        F[tickets.qr_code]
        G[payments.transaction_id]
    end
    
    style A fill:#E3F2FD
    style B fill:#E8F5E9
    style C fill:#FFF3E0
    style D fill:#FCE4EC
```

---

## 🌐 API Design

### API Structure

```mermaid
graph LR
    subgraph "API v1"
        A[/api/v1/auth]
        B[/api/v1/users]
        C[/api/v1/events]
        D[/api/v1/tickets]
        E[/api/v1/orders]
        F[/api/v1/payments]
        G[/api/v1/checkin]
        H[/api/v1/admin]
    end
    
    style A fill:#4CAF50
    style F fill:#FF9800
    style G fill:#2196F3
    style H fill:#F44336
```

### Key Endpoints

#### Authentication
```
POST   /api/v1/auth/register          - User registration
POST   /api/v1/auth/login             - User login
POST   /api/v1/auth/refresh           - Refresh token
POST   /api/v1/auth/logout            - Logout
GET    /api/v1/auth/oauth/google      - OAuth login
```

#### Events
```
GET    /api/v1/events                 - List events (paginated)
GET    /api/v1/events/:id             - Get event details
POST   /api/v1/events                 - Create event (organizer)
PUT    /api/v1/events/:id             - Update event (organizer)
DELETE /api/v1/events/:id             - Delete event (organizer)
GET    /api/v1/events/search          - Search events
```

#### Tickets
```
POST   /api/v1/tickets/reserve        - Reserve ticket (15min hold)
GET    /api/v1/tickets/:id            - Get ticket details
POST   /api/v1/tickets/:id/transfer   - Transfer ticket
GET    /api/v1/tickets/my-tickets     - User's tickets
```

#### Orders
```
POST   /api/v1/orders                 - Create order
GET    /api/v1/orders/:id             - Get order details
GET    /api/v1/orders/my-orders       - User's orders
POST   /api/v1/orders/:id/refund      - Request refund
```

#### Check-in
```
POST   /api/v1/checkin/validate       - Validate QR code
GET    /api/v1/checkin/event/:id      - Event check-in stats
```

### Response Format

```json
{
  "success": true,
  "data": {...},
  "message": "Operation successful",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Format

```json
{
  "success": false,
  "error": {
    "code": "TICKET_NOT_AVAILABLE",
    "message": "Ticket tier sold out",
    "details": {}
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 📁 Project Structure

```
ticketing-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
│
├── internal/                        # Private application code
│   ├── auth/                       # Authentication domain
│   │   ├── handler.go              # HTTP handlers
│   │   ├── service.go              # Business logic
│   │   ├── repository.go           # Data access
│   │   └── model.go                # Domain models
│   │
│   ├── users/                      # User management
│   ├── events/                     # Event management
│   ├── tickets/                    # Ticket operations
│   ├── orders/                     # Order processing
│   ├── payments/                   # Payment integration
│   ├── checkin/                    # Check-in system
│   ├── notifications/              # Notification service
│   └── admin/                      # Admin operations
│
├── pkg/                            # Shared packages
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth.go
│   │   ├── ratelimit.go
│   │   └── logger.go
│   ├── database/                   # DB connection
│   ├── cache/                      # Redis client
│   ├── logger/                     # Logging utilities
│   ├── jwt/                        # JWT utilities
│   ├── queue/                      # Message queue
│   └── utils/                      # Helper functions
│
├── configs/                        # Configuration files
│   └── config.yaml
│
├── migrations/                     # Database migrations
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   └── ...
│
├── scripts/                        # Utility scripts
│   ├── seed.go                     # Database seeding
│   └── deploy.sh                   # Deployment script
│
├── tests/                          # Test files
│   ├── integration/
│   └── unit/
│
├── .env.example                    # Environment template
├── docker-compose.yml              # Local development
├── Dockerfile                      # Production image
├── go.mod                          # Go dependencies
├── go.sum                          # Dependency checksums
└── README.md                       # This file
```

### Domain Layer Architecture

```mermaid
graph TB
    subgraph "Each Domain Module"
        Handler[Handler Layer<br/>HTTP/gRPC]
        Service[Service Layer<br/>Business Logic]
        Repository[Repository Layer<br/>Data Access]
        Model[Model Layer<br/>Domain Entities]
        
        Handler --> Service
        Service --> Repository
        Repository --> Model
    end
    
    subgraph "Dependencies"
        DB[(Database)]
        Cache[(Cache)]
        Queue[Message Queue]
    end
    
    Repository --> DB
    Service --> Cache
    Service --> Queue
    
    style Handler fill:#4CAF50
    style Service fill:#2196F3
    style Repository fill:#FF9800
    style Model fill:#9C27B0
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/ticketing-backend.git
   cd ticketing-backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start services with Docker**
   ```bash
   docker-compose up -d postgres redis kafka
   ```

5. **Run migrations**
   ```bash
   make migrate-up
   ```

6. **Seed database (optional)**
   ```bash
   go run scripts/seed.go
   ```

7. **Start the server**
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

### Environment Variables

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ticket_booking

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

# OAuth
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=

# Payments
PAYSTACK_SECRET_KEY=
STRIPE_SECRET_KEY=
WEBHOOK_SECRET=

# Kafka
KAFKA_BROKERS=localhost:9092

# S3
S3_BUCKET=
S3_REGION=
S3_ACCESS_KEY=
S3_SECRET_KEY=
```

---

## 🔒 Security

### Security Measures

```mermaid
graph TD
    A[Security Layers]
    
    A --> B[Authentication]
    B --> B1[JWT Tokens]
    B --> B2[OAuth 2.0]
    B --> B3[Refresh Tokens]
    
    A --> C[Authorization]
    C --> C1[Role-Based Access Control]
    C --> C2[Resource Ownership]
    C --> C3[Permission Checks]
    
    A --> D[Data Protection]
    D --> D1[bcrypt/argon2 Password Hash]
    D --> D2[Encrypted Tokens]
    D --> D3[PCI Compliance]
    
    A --> E[API Protection]
    E --> E1[Rate Limiting]
    E --> E2[Request Validation]
    E --> E3[Webhook Signatures]
    
    A --> F[Audit]
    F --> F1[Audit Logs]
    F --> F2[Payment Trail]
    F --> F3[Check-in History]
    
    style B fill:#4CAF50
    style C fill:#2196F3
    style D fill:#FF9800
    style E fill:#9C27B0
    style F fill:#F44336
```

### Security Checklist

- ✅ Password hashing with bcrypt/argon2
- ✅ JWT with short expiry + refresh tokens
- ✅ Rate limiting on all endpoints
- ✅ RBAC (Role-Based Access Control)
- ✅ Webhook signature verification
- ✅ SQL injection prevention (ORM)
- ✅ XSS protection
- ✅ CORS configuration
- ✅ Audit logging for sensitive operations
- ✅ PCI DSS compliance for payments

---

## ⚡ Performance

### Performance Targets

| Metric | Target | Strategy |
|--------|--------|----------|
| **API Response Time** | < 200ms | Caching, indexing, connection pooling |
| **Uptime** | 99.9% | Load balancing, health checks, auto-scaling |
| **Concurrent Users** | 10,000+ | Horizontal scaling, event-driven architecture |
| **Database Queries** | < 50ms | Proper indexing, query optimization |
| **Cache Hit Rate** | > 80% | Redis caching strategy |

### Caching Strategy

```mermaid
graph LR
    A[Request] --> B{Cache Check}
    B -->|Hit| C[Return Cached Data]
    B -->|Miss| D[Query Database]
    D --> E[Update Cache]
    E --> F[Return Data]
    
    subgraph "Cache Layers"
        G[Event List - 5min TTL]
        H[Event Details - 15min TTL]
        I[Ticket Reservation - 15min TTL]
        J[User Session - 24h TTL]
    end
    
    style C fill:#4CAF50
    style D fill:#FF9800
```

### Observability

```mermaid
graph TB
    subgraph "Monitoring"
        A[Prometheus Metrics]
        B[OpenTelemetry Traces]
        C[Structured Logs]
    end
    
    subgraph "Dashboards"
        D[Grafana]
        E[Alerting]
    end
    
    subgraph "Key Metrics"
        F[Request Rate]
        G[Error Rate]
        H[Latency P99]
        I[Database Connections]
        J[Cache Hit Rate]
    end
    
    A --> D
    B --> D
    C --> D
    D --> E
    
    D --> F
    D --> G
    D --> H
    D --> I
    D --> J
```

---

## 🎯 Success Metrics

```mermaid
graph LR
    A[Business Metrics]
    
    A --> B[Conversion Rate<br/>> 70%]
    A --> C[Payment Success<br/>> 95%]
    A --> D[Ticket Scan Success<br/>> 99%]
    A --> E[Organizer Retention<br/>> 60%]
    A --> F[Refund Frequency<br/>< 5%]
    
    style B fill:#4CAF50
    style C fill:#4CAF50
    style D fill:#4CAF50
    style E fill:#2196F3
    style F fill:#FF9800
```

---

## 📝 Future Roadmap

```mermaid
gantt
    title Product Roadmap
    dateFormat YYYY-MM-DD
    section Phase 1 - MVP
    Core API Development           :2024-01-01, 60d
    Payment Integration           :2024-02-01, 30d
    QR System                     :2024-02-15, 20d
    section Phase 2 - Enhancement
    Mobile Apps                   :2024-03-01, 45d
    Analytics Dashboard          :2024-03-15, 30d
    Push Notifications           :2024-04-01, 20d
    section Phase 3 - Scale
    Microservices Migration      :2024-05-01, 60d
    Advanced Analytics           :2024-06-01, 30d
    Multi-currency Support       :2024-06-15, 20d
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📧 Contact

For questions or support, please contact:
- **Email**: support@eventix.com
- **Website**: https://eventix.com
- **Documentation**: https://docs.eventix.com

---

<div align="center">
  <p>© 2025 Eventix. All rights reserved.</p>
</div>
