# 🚀 Setup Guide - Blog API

This guide will walk you through setting up the Blog API project from scratch.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21 or higher** - [Download here](https://golang.org/dl/)
- **MongoDB** - [Download here](https://www.mongodb.com/try/download/community)
- **Git** - [Download here](https://git-scm.com/downloads)
- **Postman** (optional) - [Download here](https://www.postman.com/downloads/)

## 🛠️ Installation Steps

### Step 1: Clone the Repository

```bash
git clone https://github.com/BemnetMussa/g6-Starter_project.git
cd g6-Starter_project
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Set Up Environment Variables

Create a `.env` file in the root directory:

```bash
# Copy the example environment file
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Server Configuration
APP_PORT=8080
APP_BASE_URL=http://localhost:8080

# Database Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=blog_api

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here-make-it-long-and-random
JWT_ACCESS_TOKEN_EXPIRY=2h
JWT_REFRESH_TOKEN_EXPIRY=7d

# Email Configuration (SMTP) - Optional
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# AI Service Configuration - Choose one
HUGGING_FACE_TOKEN=your-hugging-face-token
GROK_API_TOKEN=your-grok-api-token
OPENROUTER_API_TOKEN=your-openrouter-api-token
```

### Step 4: Set Up MongoDB

#### Option A: Local MongoDB

1. **Install MongoDB** on your system
2. **Start MongoDB service**:

   ```bash
   # Windows
   net start MongoDB

   # macOS/Linux
   sudo systemctl start mongod
   ```

#### Option B: MongoDB Atlas (Cloud)

1. **Create a MongoDB Atlas account** at [mongodb.com](https://www.mongodb.com/cloud/atlas)
2. **Create a new cluster**
3. **Get your connection string** and update `MONGODB_URI` in `.env`

### Step 5: Set Up Email (Optional)

#### For Gmail:

1. **Enable 2-Factor Authentication** on your Google account
2. **Generate an App Password**:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate a new app password
3. **Update `.env`** with your Gmail credentials

#### For Other Providers:

Update the SMTP settings in `.env` according to your email provider's specifications.

### Step 6: Set Up AI Services (Optional)

#### Hugging Face:

1. **Create account** at [huggingface.co](https://huggingface.co)
2. **Get API token** from settings
3. **Add to `.env`**:
   ```env
   HUGGING_FACE_TOKEN=your-token-here
   ```

#### OpenRouter:

1. **Create account** at [openrouter.ai](https://openrouter.ai)
2. **Get API key** from settings
3. **Add to `.env`**:
   ```env
   OPENROUTER_API_TOKEN=your-key-here
   ```

## 🚀 Running the Application

### Development Mode (with Hot Reload)

```bash
# Install Air for hot reload (if not installed)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Production Mode

```bash
go run Delivery/main.go
```

### Using Go Modules

```bash
go mod tidy
go run .
```

## 🧪 Testing the Setup

### 1. Check Server Status

Visit: `http://localhost:8080`

You should see the server running.

### 2. Test Registration

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Check Console for Verification Link

If SMTP is not configured, check the console output for the verification link.

### 4. Verify Email

Click the verification link or use curl:

```bash
curl -X GET "http://localhost:8080/auth/verify?token=YOUR_TOKEN"
```

### 5. Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## 📱 Using Postman

### 1. Import Collection

1. **Open Postman**
2. **Import** the collection from `docs/postman_collection.json`
3. **Set up environment variables**:
   - `base_url`: `http://localhost:8080`
   - `access_token`: (will be set after login)

### 2. Test Flow

1. **Register** a new user
2. **Check console** for verification link
3. **Verify email** using the link
4. **Login** to get JWT token
5. **Update environment** with the token
6. **Test other endpoints**

## 🔧 Troubleshooting

### Common Issues

#### 1. MongoDB Connection Error

**Error:** `Failed to connect to MongoDB`

**Solution:**

- Ensure MongoDB is running
- Check `MONGODB_URI` in `.env`
- For Atlas, ensure IP is whitelisted

#### 2. Port Already in Use

**Error:** `Failed to start server: listen tcp :8080: bind: address already in use`

**Solution:**

```bash
# Change port in .env
APP_PORT=8081

# Or kill the process using port 8080
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# macOS/Linux
lsof -ti:8080 | xargs kill -9
```

#### 3. Email Not Sending

**Error:** No verification emails received

**Solution:**

- Check console logs for verification links
- Verify SMTP settings in `.env`
- For Gmail, ensure app password is correct

#### 4. JWT Secret Error

**Error:** `JWT_SECRET environment variable is required`

**Solution:**

- Ensure `.env` file exists
- Check `JWT_SECRET` is set
- Restart the application

#### 5. AI Service Errors

**Error:** `API request failed`

**Solution:**

- Check API tokens in `.env`
- Verify API service is available
- Check rate limits

### Debug Mode

Enable debug logging by setting:

```env
GIN_MODE=debug
```

## 📊 Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

### Logs

Check application logs for:

- Database connections
- Email sending status
- API requests
- Error messages

## 🔒 Security Considerations

### Production Setup

1. **Use strong JWT secret**:

   ```env
   JWT_SECRET=your-very-long-and-random-secret-key-here
   ```

2. **Enable HTTPS** in production

3. **Set up proper MongoDB authentication**

4. **Use environment-specific configurations**

5. **Implement proper logging and monitoring**

### Environment Variables Best Practices

- **Never commit `.env` files** to version control
- **Use different configurations** for development and production
- **Rotate secrets regularly**
- **Use secure secret management** in production

## 🚀 Deployment

### Docker Deployment

1. **Build the image**:

   ```bash
   docker build -t blog-api .
   ```

2. **Run the container**:
   ```bash
   docker run -p 8080:8080 --env-file .env blog-api
   ```

### Cloud Deployment

#### Heroku:

1. **Create Heroku app**
2. **Set environment variables**
3. **Deploy**:
   ```bash
   heroku create your-app-name
   heroku config:set MONGODB_URI=your-mongodb-uri
   git push heroku main
   ```

#### AWS/GCP/Azure:

- Use container services (ECS, Cloud Run, App Service)
- Set up proper environment variables
- Configure MongoDB Atlas for database

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [JWT Documentation](https://jwt.io/)

## 🤝 Getting Help

If you encounter issues:

1. **Check the logs** for error messages
2. **Verify environment variables** are set correctly
3. **Test each component** individually
4. **Check the troubleshooting section** above
5. **Create an issue** on GitHub with detailed information

---

**Happy coding! 🚀**
