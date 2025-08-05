package handlers

import (
	"net/http"

	usecases "g6_starter_project/Usecases"
	"g6_starter_project/Infrastructure/services"

	"github.com/gin-gonic/gin"
)

type UserManagementHandler struct {
	userManagementUsecase *usecases.UserManagementUsecase
}

func NewUserManagementHandler(userManagementUsecase *usecases.UserManagementUsecase) *UserManagementHandler {
	return &UserManagementHandler{
		userManagementUsecase: userManagementUsecase,
	}
}

// PromoteUser promotes a user to admin role
func (h *UserManagementHandler) PromoteUser(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// Get admin user ID from context (set by auth middleware)
	adminID, exists := services.GinGetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin authentication required"})
		return
	}

	// Prevent admin from promoting themselves
	if adminID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot promote yourself"})
		return
	}

	// Promote the user
	updatedUser, err := h.userManagementUsecase.PromoteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user":    updatedUser,
	})
}

// DemoteUser demotes an admin user to regular user role
func (h *UserManagementHandler) DemoteUser(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// Get admin user ID from context (set by auth middleware)
	adminID, exists := services.GinGetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin authentication required"})
		return
	}

	// Prevent admin from demoting themselves
	if adminID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot demote yourself"})
		return
	}

	// Demote the user
	updatedUser, err := h.userManagementUsecase.DemoteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User demoted to regular user successfully",
		"user":    updatedUser,
	})
}

// GetUserByID returns a specific user by ID
func (h *UserManagementHandler) GetUserByID(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// Get user
	user, err := h.userManagementUsecase.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
} 