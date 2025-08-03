package handlers

import (
	"fmt"
	"net/http"

	"g6_starter_project/Domain/entities"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase *usecases.UserUsecase
	passwordResetUsecase *usecases.PasswordResetUsecase
}

func NewUserHandler(userUsecase *usecases.UserUsecase, passwordResetUsecase *usecases.PasswordResetUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		passwordResetUsecase: passwordResetUsecase,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received user data: %+v\n", user)
	fmt.Printf("Password length: %d\n", len(user.Password))

	createdUser, err := h.userUsecase.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

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

	loginUser, userToken, err := h.userUsecase.Login(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         loginUser,
		"accessToken":  userToken.AccessToken,
		"refreshToken": userToken.RefreshToken,
	})
}

// ForgotPassword handles password reset request
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.passwordResetUsecase.RequestPasswordReset(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Always return success for security (don't reveal if email exists)
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email address exists in our system, a password reset link has been sent.",
	})
}

// ResetPassword handles password reset with token
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.passwordResetUsecase.ResetPassword(request.Token, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been successfully reset.",
	})
}
