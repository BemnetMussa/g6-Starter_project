package routers

import (
	"g6_starter_project/Delivery/handlers"
	"g6_starter_project/Infrastructure/services"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userUsecase *usecases.UserUsecase, passwordResetUsecase *usecases.PasswordResetUsecase, userManagementUsecase *usecases.UserManagementUsecase, jwtService *services.JWTService) *gin.Engine {
	router := gin.Default()
	
	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUsecase, passwordResetUsecase)
	userManagementHandler := handlers.NewUserManagementHandler(userManagementUsecase)
	
	// Public routes (no authentication required)
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.POST("/forgot-password", userHandler.ForgotPassword)
	router.POST("/reset-password", userHandler.ResetPassword)
	
	// Protected admin routes
	adminGroup := router.Group("/admin")
	adminGroup.Use(services.GinAuthMiddleware(jwtService))
	adminGroup.Use(services.GinRoleAuthorization("admin"))
	{
		// User management routes
		adminGroup.PUT("/users/:id/promote", userManagementHandler.PromoteUser)
		adminGroup.PUT("/users/:id/demote", userManagementHandler.DemoteUser)
		adminGroup.GET("/users/:id", userManagementHandler.GetUserByID)
	}
	
	return router
}
