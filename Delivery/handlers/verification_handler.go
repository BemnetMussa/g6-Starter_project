package handlers

import (
	"net/http"

	"g6_starter_project/Domain/entities"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	verificationUsecase *usecases.VerificationUsecase
}

func NewVerificationHandler(verificationUC *usecases.VerificationUsecase) *VerificationHandler {
	return &VerificationHandler{
		verificationUsecase: verificationUC,
	}
}

// RegisterWithVerification handles user registration with email verification
func (h *VerificationHandler) RegisterWithVerification(c *gin.Context) {
	var user struct {
		FullName string `json:"full_name"`
		Username string `json:"username"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to User entity
	userEntity := &entities.User{
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	createdUser, err := h.verificationUsecase.RegisterWithVerification(userEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
		"user":    createdUser,
	})
}

// VerifyEmail handles email verification
func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification token is required"})
		return
	}

	err := h.verificationUsecase.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully! You can now log in to your account.",
	})
}

// ResendVerificationEmail handles resending verification email
func (h *VerificationHandler) ResendVerificationEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.verificationUsecase.ResendVerificationEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully. Please check your email.",
	})
} 