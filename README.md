# 🚀 Blog API - Go RESTful API with Email Verification & AI Integration

A comprehensive RESTful API built with Go, featuring email verification, AI-powered content generation, user management, and blog functionality. Built using Clean Architecture principles with MongoDB as the database.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Quick Start](#-quick-start)
- [Environment Variables](#-environment-variables)
- [API Documentation](#-api-documentation)
- [Authentication](#-authentication)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)

## ✨ Features

### 🔐 **Authentication & Security**

- **Email Verification System**: Mandatory email verification for all new registrations
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Password Reset**: Secure password reset via email
- **Role-Based Authorization**: Admin and user roles with different permissions
- **Rate Limiting**: API rate limiting to prevent abuse

### 👤 **User Management**

- **User Registration**: Single registration flow with email verification
- **Profile Management**: Users can update their profiles (bio, contact info, profile image)
- **User Roles**: Admin and user roles with different access levels
- **Account Verification**: Email verification required for login

### 📝 **Blog System**

- **CRUD Operations**: Create, read, update, delete blog posts
- **Like/Dislike System**: Interactive blog post reactions
- **Comments**: Comment system for blog posts
- **Search & Filtering**: Advanced search and filtering capabilities
- **Pagination**: Efficient pagination for large datasets

### 🤖 **AI Integration**

- **Content Generation**: AI-powered blog content generation
- **Topic Suggestions**: AI suggests blog topics based on categories
- **Content Enhancement**: AI improves existing blog content
- **Chat History**: Stores AI interaction history
- **Multiple AI Providers**: Support for Hugging Face, Grok, and OpenRouter

### 📧 **Email System**

- **SMTP Integration**: Real email sending with SMTP configuration
- **Verification Emails**: Professional verification email templates
- **Welcome Emails**: Automatic welcome emails after verification
- **Password Reset Emails**: Secure password reset via email
- **Fallback Logging**: Console logging when SMTP is not configured

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐         │
│  │   Handlers  │ │   Routers   │ │   Main.go   │         │
│  └─────────────┘ └─────────────┘ └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                    Use Cases Layer                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐         │
│  │ User UseCase│ │Blog UseCase │ │ AI UseCase  │         │
│  └─────────────┘ └─────────────┘ └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐         │
│  │   Entities  │ │ Repositories│ │  Interfaces │         │
│  └─────────────┘ └─────────────┘ └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐         │
│  │  MongoDB    │ │   Services  │ │  External   │         │
│  │ Repositories│ │  (JWT, AI)  │ │   APIs      │         │
│  └─────────────┘ └─────────────┘ └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Email**: SMTP with fallback logging
- **AI Integration**: Hugging Face, Grok, OpenRouter APIs
- **Password Hashing**: Bcrypt
- **Rate Limiting**: Custom rate limiter
- **Environment**: Air (hot reload for development)

## 📁 Project Structure

```
g6-Starter_project/
├── Delivery/                    # HTTP layer (controllers, routes)
│   ├── handlers/               # HTTP request handlers
│   │   ├── user_handler.go     # User authentication
│   │   ├── verification_handler.go # Email verification
│   │   ├── blog_handler.go     # Blog CRUD operations
│   │   ├── profile_handler.go  # User profile management
│   │   ├── ai_handler.go       # AI integration
│   │   ├── comment_handler.go  # Comment system
│   │   └── user_management_handler.go # Admin user management
│   ├── routers/                # Route definitions
│   │   └── router.go           # Main router setup
│   └── main.go                 # Application entry point
├── Domain/                     # Business logic layer
│   └── entities/               # Data models
│       ├── user.go             # User entity
│       ├── blog.go             # Blog entity
│       ├── comment.go          # Comment entity
│       ├── ai_chat.go          # AI chat entity
│       ├── blog_interaction.go # Blog interactions
│       └── token.go            # Token entity
├── Usecases/                   # Business logic
│   ├── user_usecase.go        # User authentication logic
│   ├── verification_usecase.go # Email verification logic
│   ├── blog_usecase.go        # Blog business logic
│   ├── profile_usecase.go     # Profile management logic
│   ├── ai_usecase.go          # AI integration logic
│   ├── comment_usecase.go     # Comment business logic
│   ├── password_reset_usecase.go # Password reset logic
│   ├── user_management_usecase.go # Admin user management
│   └── token_usecase.go       # Token management
├── Infrastructure/             # External dependencies
│   ├── services/              # External services
│   │   ├── jwt_service.go     # JWT token service
│   │   ├── email_service.go   # Email sending service
│   │   ├── ai_service.go      # AI API integration
│   │   ├── auth_middleware.go # Authentication middleware
│   │   ├── rate_limiter.go    # Rate limiting service
│   │   └── bcrypt_service.go  # Password hashing
│   ├── mongodb/               # Database layer
│   │   └── repositories/      # Data access layer
│   │       ├── user_repository_impl.go
│   │       ├── blog_repository_impl.go
│   │       ├── comment_repository_impl.go
│   │       └── chat_repository_impl.go
│   └── db/                    # Database connection
├── docs/                      # API documentation
├── .env                       # Environment variables
├── .air.toml                  # Air configuration
├── go.mod                     # Go modules
├── go.sum                     # Go dependencies
└── README.md                  # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- MongoDB (local or cloud)
- SMTP server (optional, for email functionality)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/BemnetMussa/g6-Starter_project.git
   cd g6-Starter_project
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**

   ```bash
   # Development mode (with hot reload)
   air

   # Production mode
   go run Delivery/main.go
   ```

5. **Access the API**
   ```
   http://localhost:8080
   ```

## 🔧 Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
APP_PORT=8080
APP_BASE_URL=http://localhost:8080

# Database Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=blog_api

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_ACCESS_TOKEN_EXPIRY=2h
JWT_REFRESH_TOKEN_EXPIRY=7d

# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# AI Service Configuration
# Choose one of the following:
HUGGING_FACE_TOKEN=your-hugging-face-token
GROK_API_TOKEN=your-grok-api-token
OPENROUTER_API_TOKEN=your-openrouter-api-token
```

## 📚 API Documentation

### Base URL

```
http://localhost:8080
```

### Authentication

Most endpoints require authentication via JWT Bearer token:

```
Authorization: Bearer <your-jwt-token>
```

### Endpoints Overview

#### 🔐 Authentication Endpoints

- `POST /register` - Register with email verification
- `POST /login` - User login
- `POST /logout` - User logout
- `POST /forgot-password` - Request password reset
- `POST /reset-password` - Reset password with token

#### ✅ Email Verification

- `GET /auth/verify` - Verify email with token
- `POST /auth/resend-verification` - Resend verification email

#### 👤 Profile Management

- `GET /profile/me` - Get user profile
- `PUT /profile/me` - Update user profile

#### 📝 Blog Endpoints

- `GET /blog` - List all blog posts (public)
- `GET /blog/:id` - Get specific blog post (public)
- `POST /blog` - Create new blog post (authenticated)
- `PUT /blog/:id` - Update blog post (authenticated)
- `DELETE /blog/:id` - Delete blog post (authenticated)
- `POST /blog/:id/like` - Like blog post (authenticated)
- `POST /blog/:id/dislike` - Dislike blog post (authenticated)

#### 💬 Comment System

- `POST /blog/:id/comments` - Add comment to blog post (authenticated)

#### 🤖 AI Integration

- `POST /ai/generate-content` - Generate blog content
- `POST /ai/suggest-topics` - Get topic suggestions
- `POST /ai/enhance-content` - Enhance existing content
- `GET /ai/chat-history` - Get AI chat history
- `DELETE /ai/chat/:id` - Delete specific chat

#### 👨‍💼 Admin Endpoints

- `PUT /admin/users/:id/promote` - Promote user to admin
- `PUT /admin/users/:id/demote` - Demote admin to user
- `GET /admin/users/:id` - Get user by ID

For detailed API documentation with request/response examples, see [docs/API.md](docs/API.md).

## 🔐 Authentication

### Registration Flow

1. User registers with email
2. Verification email sent automatically
3. User clicks verification link
4. Account becomes verified
5. User can now log in

### Login Process

1. User provides email and password
2. System checks if account is verified
3. If verified, JWT tokens are generated
4. Access token used for API calls

### Token Types

- **Access Token**: Short-lived (2 hours) for API access
- **Refresh Token**: Long-lived (7 days) for token renewal

## 🧪 Testing

### Using Postman

1. **Import the Postman collection** from `docs/postman_collection.json`
2. **Set up environment variables** in Postman
3. **Test the authentication flow**:
   - Register a new user
   - Verify email (check console for verification link)
   - Login to get JWT token
   - Use token for authenticated requests

### Manual Testing

```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Create a blog post (with JWT token)
curl -X POST http://localhost:8080/blog \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my blog post.",
    "tags": ["go", "api", "blog"]
  }'
```

## 🚀 Deployment

### Docker Deployment

1. **Build the Docker image**

   ```bash
   docker build -t blog-api .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 --env-file .env blog-api
   ```

### Cloud Deployment

The application can be deployed to:

- **Heroku**: Use the provided Procfile
- **AWS**: Deploy to EC2 or ECS
- **Google Cloud**: Deploy to Cloud Run
- **Azure**: Deploy to App Service

### Environment Setup for Production

1. **Set production environment variables**
2. **Configure MongoDB Atlas** for cloud database
3. **Set up SMTP service** for email functionality
4. **Configure AI API keys** for AI features
5. **Set up monitoring and logging**

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards
- Write tests for new features
- Update documentation for API changes
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues:

1. Check the [Issues](https://github.com/BemnetMussa/g6-Starter_project/issues) page
2. Create a new issue with detailed information
3. Contact the development team

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Air](https://github.com/cosmtrek/air) for hot reload

---

**Made with ❤️ using Go and Clean Architecture**
