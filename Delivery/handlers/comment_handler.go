package handlers

import (
	usecases "g6_starter_project/Usecases" // Make sure this import path is correct
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentUsecase usecases.ICommentUsecase
}

func NewCommentHandler(usecase usecases.ICommentUsecase) *CommentHandler {
	return &CommentHandler{commentUsecase: usecase}
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	blogIDStr := c.Param("id")
	userIDHex, exists := c.Get("userID") // From auth middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.commentUsecase.CreateComment(c.Request.Context(), blogIDStr, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}
