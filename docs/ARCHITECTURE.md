# 🏗️ Architecture Documentation

## Overview

The Blog API follows **Clean Architecture** principles, ensuring separation of concerns, testability, and maintainability. The system is built using Go with the Gin web framework and MongoDB as the database.

## 🏛️ Architecture Layers

```
┌─────────────────────────────────────────────────────────────────┐
│                        Delivery Layer                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   Handlers  │ │   Routers   │ │   Main.go   │            │
│  │ (Controllers)│ │ (Routes)    │ │ (Entry)     │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────┐
│                       Use Cases Layer                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │ User UseCase│ │Blog UseCase │ │ AI UseCase  │            │
│  │ (Business   │ │ (Business   │ │ (Business   │            │
│  │  Logic)     │ │  Logic)     │ │  Logic)     │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────┐
│                        Domain Layer                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   Entities  │ │ Repositories│ │  Interfaces │            │
│  │ (Data Models)│ │ (Contracts) │ │ (Contracts) │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │  MongoDB    │ │   Services  │ │  External   │            │
│  │ Repositories│ │  (JWT, AI)  │ │   APIs      │            │
│  │ (Database)  │ │ (External)  │ │ (Email, AI) │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
```

## 📁 Directory Structure

```
g6-Starter_project/
├── Delivery/                    # HTTP Layer
│   ├── handlers/               # HTTP Request Handlers
│   │   ├── user_handler.go     # User authentication
│   │   ├── verification_handler.go # Email verification
│   │   ├── blog_handler.go     # Blog CRUD operations
│   │   ├── profile_handler.go  # User profile management
│   │   ├── ai_handler.go       # AI integration
│   │   ├── comment_handler.go  # Comment system
│   │   └── user_management_handler.go # Admin operations
│   ├── routers/                # Route Definitions
│   │   └── router.go           # Main router setup
│   └── main.go                 # Application Entry Point
├── Domain/                     # Business Logic Layer
│   └── entities/               # Data Models
│       ├── user.go             # User entity
│       ├── blog.go             # Blog entity
│       ├── comment.go          # Comment entity
│       ├── ai_chat.go          # AI chat entity
│       ├── blog_interaction.go # Blog interactions
│       └── token.go            # Token entity
├── Usecases/                   # Business Logic
│   ├── user_usecase.go        # User authentication logic
│   ├── verification_usecase.go # Email verification logic
│   ├── blog_usecase.go        # Blog business logic
│   ├── profile_usecase.go     # Profile management logic
│   ├── ai_usecase.go          # AI integration logic
│   ├── comment_usecase.go     # Comment business logic
│   ├── password_reset_usecase.go # Password reset logic
│   ├── user_management_usecase.go # Admin user management
│   └── token_usecase.go       # Token management
└── Infrastructure/             # External Dependencies
    ├── services/              # External Services
    │   ├── jwt_service.go     # JWT token service
    │   ├── email_service.go   # Email sending service
    │   ├── ai_service.go      # AI API integration
    │   ├── auth_middleware.go # Authentication middleware
    │   ├── rate_limiter.go    # Rate limiting service
    │   └── bcrypt_service.go  # Password hashing
    ├── mongodb/               # Database Layer
    │   └── repositories/      # Data Access Layer
    │       ├── user_repository_impl.go
    │       ├── blog_repository_impl.go
    │       ├── comment_repository_impl.go
    │       └── chat_repository_impl.go
    └── db/                    # Database Connection
```

## 🔄 Data Flow

### Request Flow

```
1. HTTP Request → Router
2. Router → Handler (with middleware)
3. Handler → Use Case
4. Use Case → Repository
5. Repository → Database
6. Response flows back up the chain
```

### Authentication Flow

```
1. User Login Request
2. Handler validates input
3. Use Case checks credentials
4. JWT Service generates tokens
5. Response with tokens
```

### Email Verification Flow

```
1. User Registration
2. Use Case creates user (unverified)
3. Email Service sends verification
4. User clicks verification link
5. Use Case updates verification status
6. Welcome email sent
```

## 🏗️ Design Patterns

### 1. Clean Architecture

**Principles:**

- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Single Responsibility**: Each layer has a specific responsibility
- **Open/Closed**: Open for extension, closed for modification

**Benefits:**

- Testability
- Maintainability
- Flexibility
- Independence of frameworks

### 2. Repository Pattern

**Implementation:**

```go
// Domain layer defines interface
type UserRepository interface {
    CreateUser(user *User) (*User, error)
    GetUserByEmail(email string) (*User, error)
    // ... other methods
}

// Infrastructure layer implements interface
type UserRepositoryImpl struct {
    db *mongo.Collection
}

func (r *UserRepositoryImpl) CreateUser(user *User) (*User, error) {
    // MongoDB implementation
}
```

**Benefits:**

- Database abstraction
- Easy testing with mocks
- Multiple database support

### 3. Dependency Injection

**Implementation:**

```go
// Main.go wires dependencies
userRepository := repositories.NewUserRepository(database.Collection("users"))
userUseCase := usecases.NewUserUsecase(userRepository, tokenUseCase)
userHandler := handlers.NewUserHandler(userUseCase, passwordResetUseCase)
```

**Benefits:**

- Loose coupling
- Easy testing
- Flexible configuration

### 4. Middleware Pattern

**Implementation:**

```go
// Authentication middleware
func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT validation logic
    }
}

// Usage in router
protectedRoutes.Use(AuthMiddleware(jwtService))
```

**Benefits:**

- Cross-cutting concerns
- Reusable functionality
- Clean separation

## 🔐 Security Architecture

### Authentication & Authorization

**JWT Token Structure:**

```json
{
  "sub": "user_id",
  "role": "user|admin",
  "exp": "expiration_time",
  "iat": "issued_at"
}
```

**Token Types:**

- **Access Token**: Short-lived (2 hours) for API access
- **Refresh Token**: Long-lived (7 days) for token renewal

**Security Features:**

- Password hashing with bcrypt
- JWT token validation
- Role-based access control
- Rate limiting

### Email Security

**Verification Process:**

1. Secure token generation (32 bytes random)
2. 24-hour expiration
3. One-time use tokens
4. SMTP with TLS encryption

## 📊 Database Design

### Collections

#### Users Collection

```json
{
  "_id": "ObjectId",
  "full_name": "string (required)",
  "username": "string (unique)",
  "email": "string (unique, required)",
  "password": "string (hashed)",
  "role": "string (user/admin)",
  "is_verified": "boolean",
  "profile_image": "string (optional)",
  "bio": "string (optional)",
  "contact_info": {
    "phone": "string (optional)",
    "address": "string (optional)"
  },
  "reset_token": "string (optional)",
  "reset_token_expires_at": "datetime (optional)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### Blogs Collection

```json
{
  "_id": "ObjectId",
  "author_id": "ObjectId (ref: users)",
  "title": "string (required)",
  "content": "string (required)",
  "tags": ["string"],
  "view_count": "number",
  "likes": "number",
  "dislikes": "number",
  "comment_count": "number",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### Comments Collection

```json
{
  "_id": "ObjectId",
  "blog_id": "ObjectId (ref: blogs)",
  "author_id": "ObjectId (ref: users)",
  "content": "string (required)",
  "created_at": "datetime"
}
```

#### AI Chats Collection

```json
{
  "_id": "ObjectId",
  "user_id": "ObjectId (ref: users)",
  "request": "string (required)",
  "response": "string (required)",
  "tokens": "number",
  "created_at": "datetime"
}
```

### Indexes

**Users Collection:**

- `email` (unique)
- `username` (unique)
- `reset_token` (for password reset)

**Blogs Collection:**

- `author_id` (for user's posts)
- `created_at` (for sorting)
- `tags` (for filtering)

**Comments Collection:**

- `blog_id` (for post comments)
- `author_id` (for user comments)

## 🔄 API Design Patterns

### RESTful Design

**Resource-Based URLs:**

- `GET /blog` - List blog posts
- `POST /blog` - Create blog post
- `GET /blog/:id` - Get specific post
- `PUT /blog/:id` - Update post
- `DELETE /blog/:id` - Delete post

**Nested Resources:**

- `POST /blog/:id/comments` - Add comment to post
- `POST /blog/:id/like` - Like post
- `POST /blog/:id/dislike` - Dislike post

### Response Patterns

**Success Response:**

```json
{
  "data": {...},
  "message": "Success message"
}
```

**Error Response:**

```json
{
  "error": "Error message"
}
```

**Paginated Response:**

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

## 🚀 Performance Considerations

### Database Optimization

**Indexing Strategy:**

- Primary keys on `_id`
- Unique indexes on `email`, `username`
- Compound indexes for queries
- TTL indexes for tokens

**Query Optimization:**

- Projection to limit fields
- Aggregation for complex queries
- Pagination for large datasets

### Caching Strategy

**Current Implementation:**

- No caching (stateless design)
- Future: Redis for session management

**Potential Improvements:**

- Redis for frequently accessed data
- CDN for static content
- Browser caching for public endpoints

### Rate Limiting

**Implementation:**

- Token bucket algorithm
- Per-endpoint limits
- IP-based tracking

**Limits:**

- Authentication: 5 requests/minute
- Blog creation: 10 requests/hour
- AI endpoints: 20 requests/hour
- General: 100 requests/minute

## 🔧 Configuration Management

### Environment Variables

**Structure:**

```env
# Server
APP_PORT=8080
APP_BASE_URL=http://localhost:8080

# Database
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=blog_api

# Security
JWT_SECRET=your-secret-key
JWT_ACCESS_TOKEN_EXPIRY=2h
JWT_REFRESH_TOKEN_EXPIRY=7d

# External Services
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# AI Services
HUGGING_FACE_TOKEN=your-token
OPENROUTER_API_TOKEN=your-token
```

### Configuration Loading

**Implementation:**

```go
func LoadEnvVariables() {
    if err := godotenv.Load(".env"); err != nil {
        log.Println("Warning: .env file not found")
    }
}
```

## 🧪 Testing Strategy

### Testing Layers

**Unit Tests:**

- Use cases
- Services
- Repositories

**Integration Tests:**

- API endpoints
- Database operations
- External service integration

**Test Structure:**

```
tests/
├── unit/
│   ├── usecases/
│   ├── services/
│   └── repositories/
├── integration/
│   ├── handlers/
│   └── api/
└── e2e/
    └── scenarios/
```

### Mocking Strategy

**Repository Mocks:**

```go
type MockUserRepository struct {
    users map[string]*entities.User
}

func (m *MockUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
    // Mock implementation
}
```

## 📈 Scalability Considerations

### Horizontal Scaling

**Stateless Design:**

- No session storage
- JWT-based authentication
- Database as single source of truth

**Load Balancing:**

- Multiple application instances
- Database read replicas
- CDN for static content

### Vertical Scaling

**Database Optimization:**

- Proper indexing
- Query optimization
- Connection pooling

**Application Optimization:**

- Goroutine management
- Memory usage monitoring
- CPU profiling

## 🔄 Deployment Architecture

### Development Environment

**Local Setup:**

- Local MongoDB
- Hot reload with Air
- Environment variables in `.env`

### Production Environment

**Container Deployment:**

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
```

**Cloud Deployment:**

- Container orchestration (Kubernetes)
- Auto-scaling
- Load balancing
- Monitoring and logging

## 🔍 Monitoring & Logging

### Logging Strategy

**Structured Logging:**

```go
log.Printf("User %s registered successfully", user.Email)
```

**Log Levels:**

- DEBUG: Development details
- INFO: General information
- WARN: Warning messages
- ERROR: Error conditions

### Monitoring

**Health Checks:**

- Database connectivity
- External service availability
- Application status

**Metrics:**

- Request/response times
- Error rates
- Database performance
- Memory usage

## 🔒 Security Architecture

### Input Validation

**Request Validation:**

```go
type RegisterRequest struct {
    FullName string `json:"full_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### Data Protection

**Password Security:**

- Bcrypt hashing (cost factor 10)
- Salt generation
- Secure comparison

**Token Security:**

- Cryptographically secure generation
- Short expiration times
- Secure storage

### API Security

**Rate Limiting:**

- Per-endpoint limits
- IP-based tracking
- Token bucket algorithm

**CORS Configuration:**

- Origin restrictions
- Method restrictions
- Header restrictions

## 🚀 Future Enhancements

### Planned Features

**Performance:**

- Redis caching
- Database connection pooling
- Query optimization

**Security:**

- OAuth2 integration
- Two-factor authentication
- API key management

**Scalability:**

- Microservices architecture
- Event-driven design
- Message queues

**Monitoring:**

- Prometheus metrics
- Grafana dashboards
- Distributed tracing

---

This architecture provides a solid foundation for a scalable, maintainable, and secure blog API system.
