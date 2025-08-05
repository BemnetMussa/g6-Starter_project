package handlers

import (
	"net/http"

	"g6_starter_project/Domain/entities"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	userProfileUsecase *usecases.UserProfileUsecase
}

func NewUserProfileHandler(userProfileUsecase *usecases.UserProfileUsecase) *UserProfileHandler {
	return &UserProfileHandler{
		userProfileUsecase: userProfileUsecase,
	}
}

// GetMyProfile gets the current user's own profile
func (h *UserProfileHandler) GetMyProfile(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userProfileUsecase.GetUserProfileByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateMyProfile updates the current user's own profile
func (h *UserProfileHandler) UpdateMyProfile(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var updateData entities.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update the profile
	updatedUser, err := h.userProfileUsecase.UpdateUserProfile(userID.(string), &updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}
