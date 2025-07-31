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
}

func NewUserHandler(userUsecase *usecases.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
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
