package handlers

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"net/http"
	"strconv"
	"time"
	"strings"

	"g6_starter_project/Domain/entities"
	usecases "g6_starter_project/Usecases"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
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

	authorIDHex, exists := c.Get("userID")
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
	userIDHex, exists := c.Get("userID")
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

	userIDHex, _ := c.Get("userID")
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
	ctx := c.Request.Context()

	// Query params
	tag := c.Query("tag")
	authorName := c.Query("author")
	title := c.Query("title")
	sortBy := c.DefaultQuery("sortBy", "createdAt")

	// Date filtering
	var startTimePtr, endTimePtr *time.Time

	startDateStr := c.Query("startDate")
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startTimePtr = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use YYYY-MM-DD"})
			return
		}
	}

	endDateStr := c.Query("endDate")
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endTimePtr = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use YYYY-MM-DD"})
			return
		}
	}

	// Pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Usecase call
	posts, total, err := h.blogUsecase.ListPosts(ctx, tag, authorName, title, sortBy, startTimePtr, endTimePtr, int64(page), int64(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"limit": limit,
		"posts": posts,
	})
}


// DeletePost handles DELETE /posts/:id requests.
func (h *BlogHandler) DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// Get the authenticated user's info from the context.
	userIDHex, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")
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

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// LikePost handles POST /posts/:id/like requests.
func (h *BlogHandler) LikePost(c *gin.Context) {
	postID := c.Param("id")

	userIDHex, exists := c.Get("userID")
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

	userIDHex, exists := c.Get("userID")
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
