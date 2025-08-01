package handlers

import (
	"net/http"
	"starter_project/Domain/entities"
	usecases "starter_project/Usecases"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogHandler - the controller for blog-related HTTP requests.
type BlogHandler struct {
	blogUsecase usecases.IBlogUsecase
}

// NewBlogHandler is the constructor.
func NewBlogHandler(usecase usecases.IBlogUsecase) *BlogHandler {
	return &BlogHandler{blogUsecase: usecase}
}

// CreatePost handles POST /posts requests.
func (h *BlogHandler) CreatePost(c *gin.Context) {
	var post entities.Blog
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authorIDHex, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	authorID, _ := primitive.ObjectIDFromHex(authorIDHex.(string))
	post.AuthorID = authorID

	createdPost, err := h.blogUsecase.CreatePost(c.Request.Context(), &post, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, createdPost)
}

// GetPostByID handles GET /posts/:id requests.
func (h *BlogHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("id")

	// Check if a user is logged in to track the view.(optional)
	var userID *primitive.ObjectID
	userIDHex, exists := c.Get("user_id")
	if exists {
		id, err := primitive.ObjectIDFromHex(userIDHex.(string))
		if err == nil {
			userID = &id
		}
	}

	post, err := h.blogUsecase.GetPostByID(c.Request.Context(), postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost handles PUT /posts/:id requests.
func (h *BlogHandler) UpdatePost(c *gin.Context) {
	postID := c.Param("id")

	var updateData entities.Blog
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userIDHex, _ := c.Get("user_id")
	requestingUserID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	updatedPost, err := h.blogUsecase.UpdatePost(c.Request.Context(), postID, &updateData, requestingUserID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

// ListPosts handles GET /posts requests with filtering, searching, and pagination.
func (h *BlogHandler) ListPosts(c *gin.Context) {
	// Get the optional query parameters from the URL.
	tag := c.Query("tag")
	authorName := c.Query("author")
	title := c.Query("title")
	sortBy := c.Query("sort_by")

	var startDate, endDate *time.Time
	if dateStr := c.Query("start_date"); dateStr != "" {
		if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
			startDate = &t
		}
	}
	if dateStr := c.Query("end_date"); dateStr != "" {
		if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
			endDate = &t
		}
	}
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		// If 'page' is not a valid number or is less than 1, default to 1.
		page = 1
	}

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		// If 'limit' is not a valid number or is less than 1, default to 10.
		limit = 10
	}

	posts, total, err := h.blogUsecase.ListPosts(c.Request.Context(), tag, authorName, title, sortBy, startDate, endDate, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	//    the client know how many pages there are in total.
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	c.JSON(http.StatusOK, posts)
}

// DeletePost handles DELETE /posts/:id requests.
func (h *BlogHandler) DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// Get the authenticated user's info from the context.
	userIDHex, _ := c.Get("user_id")
	userRole, _ := c.Get("role")
	requestingUserID, _ := primitive.ObjectIDFromHex(userIDHex.(string))
	requestingUserRole := userRole.(string)

	err := h.blogUsecase.DeletePost(c.Request.Context(), postID, requestingUserID, requestingUserRole)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// LikePost handles POST /posts/:id/like requests.
func (h *BlogHandler) LikePost(c *gin.Context) {
	postID := c.Param("id")

	userIDHex, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	err := h.blogUsecase.LikePost(c.Request.Context(), postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

// DislikePost handles POST /posts/:id/dislike requests.
func (h *BlogHandler) DislikePost(c *gin.Context) {
	postID := c.Param("id")

	userIDHex, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	err := h.blogUsecase.DislikePost(c.Request.Context(), postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post disliked successfully"})
}
