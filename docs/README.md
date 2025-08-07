# 📚 Documentation Index

Welcome to the Blog API documentation! This comprehensive guide covers everything you need to know about the project.

## 📋 Quick Navigation

### 🚀 Getting Started

- **[Setup Guide](SETUP.md)** - Complete installation and configuration instructions
- **[API Documentation](API.md)** - Comprehensive API reference with examples
- **[Architecture Guide](ARCHITECTURE.md)** - System design and technical details

### 📖 Documentation Files

| File                                                   | Description                             | Audience                  |
| ------------------------------------------------------ | --------------------------------------- | ------------------------- |
| **[SETUP.md](SETUP.md)**                               | Step-by-step installation guide         | Developers, DevOps        |
| **[API.md](API.md)**                                   | Complete API reference with examples    | Developers, API Consumers |
| **[ARCHITECTURE.md](ARCHITECTURE.md)**                 | System architecture and design patterns | Developers, Architects    |
| **[postman_collection.json](postman_collection.json)** | Postman collection for testing          | Developers, Testers       |

## 🎯 What You'll Find Here

### For New Users

1. **Start with [SETUP.md](SETUP.md)** - Get the project running quickly
2. **Use [API.md](API.md)** - Learn how to use the API endpoints
3. **Import the Postman collection** - Test the API easily

### For Developers

1. **Read [ARCHITECTURE.md](ARCHITECTURE.md)** - Understand the system design
2. **Review [API.md](API.md)** - See all available endpoints
3. **Follow [SETUP.md](SETUP.md)** - Set up your development environment

### For API Consumers

1. **Focus on [API.md](API.md)** - Complete API reference
2. **Use the Postman collection** - Ready-to-use API tests
3. **Check authentication examples** - Learn how to authenticate

## 🚀 Quick Start

### 1. Setup (5 minutes)

```bash
git clone https://github.com/BemnetMussa/g6-Starter_project.git
cd g6-Starter_project
# Follow SETUP.md for detailed instructions
```

### 2. Test the API

```bash
# Register a user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Import Postman Collection

- Open Postman
- Import `docs/postman_collection.json`
- Set environment variables
- Start testing!

## 📚 Documentation Structure

```
docs/
├── README.md                    # This file - Documentation index
├── SETUP.md                     # Installation and setup guide
├── API.md                       # Complete API documentation
├── ARCHITECTURE.md              # System architecture guide
└── postman_collection.json      # Postman collection for testing
```

## 🔍 What's Covered

### ✅ Complete Features

- **User Authentication** - Registration, login, email verification
- **Blog Management** - CRUD operations, likes, comments
- **Profile Management** - User profiles, contact info
- **AI Integration** - Content generation, topic suggestions
- **Admin Features** - User management, role management
- **Email System** - Verification, password reset, welcome emails

### ✅ Technical Details

- **Clean Architecture** - Separation of concerns
- **RESTful API Design** - Standard HTTP methods
- **JWT Authentication** - Secure token-based auth
- **MongoDB Integration** - NoSQL database
- **Rate Limiting** - API protection
- **Error Handling** - Consistent error responses

### ✅ Development Tools

- **Postman Collection** - Ready-to-use API tests
- **Environment Configuration** - Easy setup
- **Hot Reload** - Development efficiency
- **Comprehensive Logging** - Debug information

## 🎯 Use Cases

### For Frontend Developers

- Use [API.md](API.md) to understand available endpoints
- Import Postman collection for testing
- Follow authentication flow in documentation

### For Backend Developers

- Study [ARCHITECTURE.md](ARCHITECTURE.md) for system design
- Review [SETUP.md](SETUP.md) for development environment
- Use API documentation for integration

### For DevOps Engineers

- Follow [SETUP.md](SETUP.md) for deployment
- Review architecture for scaling considerations
- Use environment configuration examples

### For API Consumers

- Focus on [API.md](API.md) for endpoint details
- Use Postman collection for testing
- Follow authentication examples

## 🔧 Common Tasks

### Setting Up Development Environment

1. **Clone repository**
2. **Install dependencies** (`go mod download`)
3. **Configure environment** (copy `.env.example` to `.env`)
4. **Start MongoDB**
5. **Run application** (`air` for development)

### Testing the API

1. **Register a user** (POST `/register`)
2. **Verify email** (check console for link)
3. **Login** (POST `/login`)
4. **Create blog post** (POST `/blog`)
5. **Test AI features** (POST `/ai/generate-content`)

### Deploying to Production

1. **Set production environment variables**
2. **Configure MongoDB Atlas**
3. **Set up SMTP for emails**
4. **Configure AI API keys**
5. **Deploy using Docker or cloud platform**

## 📞 Getting Help

### Documentation Issues

- Check if your question is answered in the docs
- Look for similar issues in the troubleshooting sections
- Review the architecture documentation for technical details

### Setup Problems

- Follow [SETUP.md](SETUP.md) step by step
- Check the troubleshooting section
- Verify environment variables are set correctly

### API Questions

- Review [API.md](API.md) for endpoint details
- Use the Postman collection for testing
- Check request/response examples

### Technical Issues

- Review [ARCHITECTURE.md](ARCHITECTURE.md) for system design
- Check logs for error messages
- Verify database connectivity

## 🚀 Next Steps

### For New Users

1. **Complete the setup** following [SETUP.md](SETUP.md)
2. **Test basic functionality** using the API documentation
3. **Explore advanced features** like AI integration
4. **Customize for your needs** based on the architecture

### For Developers

1. **Understand the architecture** from [ARCHITECTURE.md](ARCHITECTURE.md)
2. **Set up development environment** following [SETUP.md](SETUP.md)
3. **Review the codebase** with the documentation as reference
4. **Contribute improvements** based on the established patterns

### For API Consumers

1. **Import the Postman collection** for easy testing
2. **Review authentication flow** in [API.md](API.md)
3. **Test all endpoints** using the provided examples
4. **Integrate with your application** using the API reference

---

**Happy coding! 🚀**

_This documentation is maintained alongside the codebase. For the latest updates, check the repository._
