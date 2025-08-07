package handlers

import (
	"net/http"

	// "g6_starter_project/Domain/entities"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiUsecase usecases.AIUsecase
}

func NewAIHandler(aiUsecase usecases.AIUsecase) *AIHandler {
	return &AIHandler{
		aiUsecase: aiUsecase,
	}
}

// GenerateBlogContent generates blog content using AI
func (h *AIHandler) GenerateBlogContent(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var request struct {
		Topic string `json:"topic" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Generate content
	chat, err := h.aiUsecase.GenerateBlogContent(userID.(string), request.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// SuggestTopics generates topic suggestions using AI
func (h *AIHandler) SuggestTopics(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var request struct {
		Category string `json:"category" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Generate topics
	chat, err := h.aiUsecase.SuggestTopics(userID.(string), request.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// EnhanceContent improves existing content using AI
func (h *AIHandler) EnhanceContent(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var request struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Enhance content
	chat, err := h.aiUsecase.EnhanceContent(userID.(string), request.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// GetChatHistory gets user's AI chat history
func (h *AIHandler) GetChatHistory(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get chat history
	chats, err := h.aiUsecase.GetChatHistory(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// DeleteChat deletes a specific chat
func (h *AIHandler) DeleteChat(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get chat ID from URL parameter
	chatID := c.Param("id")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat ID is required"})
		return
	}

	// Delete chat
	err := h.aiUsecase.DeleteChat(userID.(string), chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
} 