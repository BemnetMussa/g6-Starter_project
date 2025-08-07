# 📚 API Documentation

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Error Responses](#error-responses)
- [Authentication Endpoints](#authentication-endpoints)
- [Email Verification Endpoints](#email-verification-endpoints)
- [Profile Management Endpoints](#profile-management-endpoints)
- [Blog Endpoints](#blog-endpoints)
- [Comment Endpoints](#comment-endpoints)
- [AI Integration Endpoints](#ai-integration-endpoints)
- [Admin Endpoints](#admin-endpoints)

## Base URL

```
http://localhost:8080
```

## Authentication

Most endpoints require authentication via JWT Bearer token:

```
Authorization: Bearer <your-jwt-token>
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

---

## Authentication Endpoints

### 1. Register User

**Endpoint:** `POST /register`

**Description:** Register a new user with email verification

**Request Body:**

```json
{
  "full_name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**

```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "user": {
    "id": "68948f61ac1badb0de2ac59c",
    "full_name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "password": "",
    "is_verified": false,
    "created_at": "2025-08-07T11:35:34.440Z",
    "updated_at": "2025-08-07T11:35:34.440Z"
  }
}
```

**Notes:**

- Email verification is required before login
- Verification email is sent automatically
- Check console logs for verification link if SMTP is not configured

---

### 2. Login User

**Endpoint:** `POST /login`

**Description:** Authenticate user and get JWT tokens

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "user": {
    "id": "68948f61ac1badb0de2ac59c",
    "full_name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "password": "",
    "role": "user",
    "is_verified": true,
    "profile_image": null,
    "bio": null,
    "contact_info": null,
    "created_at": "2025-08-07T11:35:34.440Z",
    "updated_at": "2025-08-07T11:35:34.440Z"
  },
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (401 Unauthorized):**

```json
{
  "error": "account not verified. Please check your email and verify your account"
}
```

---

### 3. Logout User

**Endpoint:** `POST /logout`

**Description:** Logout user and invalidate tokens

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "logged out successfully"
}
```

---

### 4. Forgot Password

**Endpoint:** `POST /forgot-password`

**Description:** Request password reset email

**Request Body:**

```json
{
  "email": "john@example.com"
}
```

**Response (200 OK):**

```json
{
  "message": "If the email exists, a password reset link has been sent."
}
```

---

### 5. Reset Password

**Endpoint:** `POST /reset-password`

**Description:** Reset password using reset token

**Request Body:**

```json
{
  "token": "reset-token-from-email",
  "new_password": "newpassword123"
}
```

**Response (200 OK):**

```json
{
  "message": "Password reset successfully."
}
```

---

## Email Verification Endpoints

### 1. Verify Email

**Endpoint:** `GET /auth/verify`

**Description:** Verify email with verification token

**Query Parameters:**

- `token` (required): Verification token from email

**URL Example:**

```
GET /auth/verify?token=2e0274f12b97aa9667b54af7f59cfe9cb1086f4841f69d32ca2c8cff86dad741
```

**Response (200 OK):**

```json
{
  "message": "Email verified successfully! You can now log in to your account."
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "verification token is required"
}
```

---

### 2. Resend Verification Email

**Endpoint:** `POST /auth/resend-verification`

**Description:** Resend verification email to user

**Request Body:**

```json
{
  "email": "john@example.com"
}
```

**Response (200 OK):**

```json
{
  "message": "Verification email sent successfully. Please check your email."
}
```

---

## Profile Management Endpoints

### 1. Get My Profile

**Endpoint:** `GET /profile/me`

**Description:** Get current user's profile

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "id": "68948f61ac1badb0de2ac59c",
  "full_name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user",
  "is_verified": true,
  "profile_image": "https://example.com/profile.jpg",
  "bio": "Software developer passionate about Go and clean architecture",
  "contact_info": {
    "phone": "+1234567890",
    "address": "123 Main St, City, Country"
  },
  "created_at": "2025-08-07T11:35:34.440Z",
  "updated_at": "2025-08-07T11:35:34.440Z"
}
```

---

### 2. Update My Profile

**Endpoint:** `PUT /profile/me`

**Description:** Update current user's profile

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "full_name": "John Doe Updated",
  "username": "johndoe_updated",
  "profile_image": "https://example.com/new-profile.jpg",
  "bio": "This is my updated bio!",
  "contact_info": {
    "phone": "+1234567890",
    "address": "123 Updated Street, City"
  }
}
```

**Response (200 OK):**

```json
{
  "id": "68948f61ac1badb0de2ac59c",
  "full_name": "John Doe Updated",
  "username": "johndoe_updated",
  "email": "john@example.com",
  "role": "user",
  "is_verified": true,
  "profile_image": "https://example.com/new-profile.jpg",
  "bio": "This is my updated bio!",
  "contact_info": {
    "phone": "+1234567890",
    "address": "123 Updated Street, City"
  },
  "created_at": "2025-08-07T11:35:34.440Z",
  "updated_at": "2025-08-07T11:59:41.453Z"
}
```

---

## Blog Endpoints

### 1. List Blog Posts

**Endpoint:** `GET /blog`

**Description:** Get all blog posts with pagination and filtering

**Query Parameters:**

- `title` (optional): Search by title (case-insensitive)
- `author` (optional): Filter by author name
- `tag` (optional): Filter by tags (comma-separated)
- `startDate` (optional): Filter by start date (YYYY-MM-DD)
- `endDate` (optional): Filter by end date (YYYY-MM-DD)
- `sortBy` (optional): Sort order (popularity, date_asc, date_desc)
- `minPopularity` (optional): Minimum likes count
- `maxPopularity` (optional): Maximum likes count
- `page` (optional): Page number (default: 1)
- `limit` (optional): Posts per page (default: 10)

**URL Example:**

```
GET /blog?page=1&limit=10&sortBy=date_desc&tag=go,api
```

**Response (200 OK):**

```json
{
  "limit": 10,
  "page": 1,
  "posts": [
    {
      "id": "689457b56e2cae04a9ace74d",
      "author_id": "6893544d594f56c731efd47d",
      "title": "Psychology",
      "content": "Psychology is the scientific study of the mind and behavior...",
      "tags": ["mind", "science", "psychology"],
      "view_count": 0,
      "likes": 0,
      "dislikes": 0,
      "comment_count": 0,
      "created_at": "2025-08-07T07:37:25.509Z",
      "updated_at": "2025-08-07T07:37:25.509Z"
    }
  ],
  "total": 2
}
```

---

### 2. Get Blog Post by ID

**Endpoint:** `GET /blog/:id`

**Description:** Get specific blog post by ID

**URL Example:**

```
GET /blog/689457b56e2cae04a9ace74d
```

**Response (200 OK):**

```json
{
  "id": "689457b56e2cae04a9ace74d",
  "author_id": "6893544d594f56c731efd47d",
  "title": "Psychology",
  "content": "Psychology is the scientific study of the mind and behavior...",
  "tags": ["mind", "science", "psychology"],
  "view_count": 0,
  "likes": 0,
  "dislikes": 0,
  "comment_count": 0,
  "created_at": "2025-08-07T07:37:25.509Z",
  "updated_at": "2025-08-07T07:37:25.509Z"
}
```

---

### 3. Create Blog Post

**Endpoint:** `POST /blog`

**Description:** Create a new blog post

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "title": "Psychology",
  "content": "Psychology is the scientific study of the mind and behavior, exploring how people think, feel, and act. It helps us understand mental processes, emotions, and social interactions.",
  "tags": ["mind", "science", "psychology"]
}
```

**Response (201 Created):**

```json
{
  "id": "689457b56e2cae04a9ace74d",
  "author_id": "6893544d594f56c731efd47d",
  "title": "Psychology",
  "content": "Psychology is the scientific study of the mind and behavior...",
  "tags": ["mind", "science", "psychology"],
  "view_count": 0,
  "likes": 0,
  "dislikes": 0,
  "comment_count": 0,
  "created_at": "2025-08-07T07:37:25.509Z",
  "updated_at": "2025-08-07T07:37:25.509Z"
}
```

---

### 4. Update Blog Post

**Endpoint:** `PUT /blog/:id`

**Description:** Update an existing blog post (author only)

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "title": "Programming Language",
  "content": "A programming language is a formal set of instructions used to communicate with computers and create software applications...",
  "tags": ["Go", "Python", "Java"]
}
```

**Response (200 OK):**

```json
{
  "id": "68936643594f56c731efd482",
  "author_id": "68935bee594f56c731efd47f",
  "title": "Programming Language",
  "content": "A programming language is a formal set of instructions...",
  "tags": ["Go", "Python", "Java"],
  "view_count": 0,
  "likes": 0,
  "dislikes": 0,
  "comment_count": 2,
  "created_at": "2025-08-06T14:27:15.66Z",
  "updated_at": "2025-08-07T00:41:30.412Z"
}
```

---

### 5. Delete Blog Post

**Endpoint:** `DELETE /blog/:id`

**Description:** Delete a blog post (author only)

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "Post deleted successfully"
}
```

---

### 6. Like Blog Post

**Endpoint:** `POST /blog/:id/like`

**Description:** Like a blog post

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "Post liked successfully"
}
```

---

### 7. Dislike Blog Post

**Endpoint:** `POST /blog/:id/dislike`

**Description:** Dislike a blog post

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "Post disliked successfully"
}
```

---

## Comment Endpoints

### 1. Create Comment

**Endpoint:** `POST /blog/:id/comments`

**Description:** Add a comment to a blog post

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "content": "This is the first comment from the user"
}
```

**Response (201 Created):**

```json
{
  "id": "68945ac66e2cae04a9ace74e",
  "blog_id": "689457b56e2cae04a9ace74d",
  "author_id": "68935bee594f56c731efd47f",
  "content": "This is the first comment from the user",
  "created_at": "2025-08-07T00:50:30.896Z"
}
```

---

## AI Integration Endpoints

### 1. Generate Blog Content

**Endpoint:** `POST /ai/generate-content`

**Description:** Generate blog content using AI

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "topic": "Go programming best practices"
}
```

**Response (200 OK):**

```json
{
  "chat": {
    "id": "68930aecca52e72af4e41c60",
    "user_id": "688f21291e97242755bd7393",
    "request": "Go programming best practices",
    "response": "Title: Mastering Go Programming: Best Practices for Efficient and Clean Code\n\nIntroduction:\nGo, also known as Golang, is a powerful programming language that has gained popularity for its simplicity, efficiency, and performance...",
    "tokens": 682,
    "created_at": "2025-08-06T07:57:32.124Z"
  }
}
```

---

### 2. Suggest Topics

**Endpoint:** `POST /ai/suggest-topics`

**Description:** Get AI-generated topic suggestions

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "category": "Web Development"
}
```

**Response (200 OK):**

```json
{
  "chat": {
    "id": "68930a80ca52e72af4e41c5f",
    "user_id": "688f21291e97242755bd7393",
    "request": "Web Development",
    "response": "1. \"The Evolution of Web Development: From HTML to AI-Powered Websites\"\n2. \"Mastering Responsive Design: Creating Websites for Every Screen Size\"\n3. \"The Power of JavaScript Frameworks: A Guide to Choosing the Right One\"...",
    "tokens": 327,
    "created_at": "2025-08-06T07:55:44.263Z"
  }
}
```

---

### 3. Enhance Content

**Endpoint:** `POST /ai/enhance-content`

**Description:** Enhance existing blog content using AI

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "content": "This is a basic blog post about Go programming. It needs improvement."
}
```

**Response (200 OK):**

```json
{
  "chat": {
    "id": "68930aecca52e72af4e41c60",
    "user_id": "688f21291e97242755bd7393",
    "request": "This is a basic blog post about Go programming. It needs improvement.",
    "response": "Title: Mastering Go Programming: A Comprehensive Guide for Beginners\n\nIntroduction:\nWelcome to our comprehensive guide on Go programming!...",
    "tokens": 453,
    "created_at": "2025-08-06T07:57:32.124Z"
  }
}
```

---

### 4. Get Chat History

**Endpoint:** `GET /ai/chat-history`

**Description:** Get user's AI chat history

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "chats": [
    {
      "id": "6892fe07c4b05d8698d5b0de",
      "user_id": "688f21291e97242755bd7393",
      "request": "Go programming best practices",
      "response": "This is a mock AI response for the prompt...",
      "tokens": 85,
      "created_at": "2025-08-06T07:02:31.999Z"
    },
    {
      "id": "68930aecca52e72af4e41c60",
      "user_id": "688f21291e97242755bd7393",
      "request": "Go programming best practices",
      "response": "Title: Mastering Go Programming: Best Practices...",
      "tokens": 682,
      "created_at": "2025-08-06T07:51:40.77Z"
    }
  ]
}
```

---

### 5. Delete Chat

**Endpoint:** `DELETE /ai/chat/:id`

**Description:** Delete a specific AI chat

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "Chat deleted successfully"
}
```

---

## Admin Endpoints

### 1. Promote User

**Endpoint:** `PUT /admin/users/:id/promote`

**Description:** Promote a user to admin role

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "User promoted to admin successfully",
  "user": {
    "id": "68948f61ac1badb0de2ac59c",
    "full_name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "admin",
    "is_verified": true,
    "created_at": "2025-08-07T11:35:34.440Z",
    "updated_at": "2025-08-07T11:59:41.453Z"
  }
}
```

---

### 2. Demote User

**Endpoint:** `PUT /admin/users/:id/demote`

**Description:** Demote an admin to user role

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "User demoted to user successfully",
  "user": {
    "id": "68948f61ac1badb0de2ac59c",
    "full_name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user",
    "is_verified": true,
    "created_at": "2025-08-07T11:35:34.440Z",
    "updated_at": "2025-08-07T11:59:41.453Z"
  }
}
```

---

### 3. Get User by ID

**Endpoint:** `GET /admin/users/:id`

**Description:** Get user information by ID

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**

```json
{
  "id": "68948f61ac1badb0de2ac59c",
  "full_name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user",
  "is_verified": true,
  "profile_image": null,
  "bio": null,
  "contact_info": null,
  "created_at": "2025-08-07T11:35:34.440Z",
  "updated_at": "2025-08-07T11:59:41.453Z"
}
```

---

## Data Models

### User Entity

```json
{
  "id": "string (ObjectID)",
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

### Blog Entity

```json
{
  "id": "string (ObjectID)",
  "author_id": "string (ObjectID)",
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

### Comment Entity

```json
{
  "id": "string (ObjectID)",
  "blog_id": "string (ObjectID)",
  "author_id": "string (ObjectID)",
  "content": "string (required)",
  "created_at": "datetime"
}
```

### AI Chat Entity

```json
{
  "id": "string (ObjectID)",
  "user_id": "string (ObjectID)",
  "request": "string (required)",
  "response": "string (required)",
  "tokens": "number",
  "created_at": "datetime"
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Authentication endpoints**: 5 requests per minute
- **Blog creation**: 10 requests per hour
- **AI endpoints**: 20 requests per hour
- **General endpoints**: 100 requests per minute

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

---

## Testing Examples

### Complete User Flow

1. **Register a new user:**

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

2. **Verify email (check console for link):**

```bash
curl -X GET "http://localhost:8080/auth/verify?token=YOUR_TOKEN"
```

3. **Login to get JWT token:**

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

4. **Create a blog post:**

```bash
curl -X POST http://localhost:8080/blog \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my blog post.",
    "tags": ["go", "api", "blog"]
  }'
```

5. **Generate AI content:**

```bash
curl -X POST http://localhost:8080/ai/generate-content \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "Go programming best practices"
  }'
```

---

## Notes

- All timestamps are in ISO 8601 format (UTC)
- ObjectIDs are MongoDB ObjectIDs
- JWT tokens expire after 2 hours (access) and 7 days (refresh)
- Email verification tokens expire after 24 hours
- Password reset tokens expire after 15 minutes
- AI responses are stored in the database for history tracking
- Rate limiting is applied to prevent API abuse
- All sensitive data (passwords, tokens) are hashed/encrypted
