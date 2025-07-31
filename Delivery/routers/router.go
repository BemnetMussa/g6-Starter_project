package routers

import (
	"github.com/gin-gonic/gin"
	"g6-Starter_project/Delivery/handlers"
	usecases "g6-Starter_project/Usecases"
)

func SetupRouter(userUsecase *usecases.UserUsecase) *gin.Engine {
	router := gin.Default()
	userHandler := handlers.NewUserHandler(userUsecase)
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	return router
}	

