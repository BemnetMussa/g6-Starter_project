package routers

import (
	"g6_starter_project/Delivery/handlers"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userUsecase *usecases.UserUsecase) *gin.Engine {
	router := gin.Default()
	userHandler := handlers.NewUserHandler(userUsecase)
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	return router
}
