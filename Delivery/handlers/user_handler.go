package handlers

import (
	"fmt"
	"net/http"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/services"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase           *usecases.UserUsecase
	passwordResetUsecase  *usecases.PasswordResetUsecase
}

// NewUserHandler initializes a new UserHandler
func NewUserHandler(userUC *usecases.UserUsecase, resetUC *usecases.PasswordResetUsecase) *UserHandler {
	return &UserHandler{
		userUsecase:          userUC,
		passwordResetUsecase: resetUC,
	}
}

// Login handles user authentication and token generation
func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := entities.User{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	authenticatedUser, token, err := h.userUsecase.Login(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         authenticatedUser,
		"accessToken":  token.AccessToken,
		"refreshToken": token.RefreshToken,
	})
}

// ForgotPassword sends a reset link to the user's email
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.passwordResetUsecase.RequestPasswordReset(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't expose whether email exists or not
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent.",
	})
}

// ResetPassword sets a new password for the user using a reset token
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.passwordResetUsecase.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully."})
}


func (h *UserHandler) Logout(c *gin.Context) {
	fmt.Println("logout handler is processing ")
	userID, exists := services.GinGetUserID(c)
	fmt.Println("userId output: ", userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	err := h.userUsecase.Logout(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
