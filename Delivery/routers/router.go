package routers

import (
	"g6_starter_project/Infrastructure/services"
	"g6_starter_project/Delivery/handlers"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userUsecase *usecases.UserUsecase,
	passwordResetUsecase *usecases.PasswordResetUsecase,
	userManagementUsecase *usecases.UserManagementUsecase,
	blogHandler *handlers.BlogHandler,
	jwtService *services.JWTService,
) *gin.Engine {

	router := gin.Default()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUsecase, passwordResetUsecase)
	userManagementHandler := handlers.NewUserManagementHandler(userManagementUsecase)

	// Public routes
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.POST("/forgot-password", userHandler.ForgotPassword)
	router.POST("/reset-password", userHandler.ResetPassword)

	// Blog routes
	postRoutes := router.Group("/blog")
	{
		// Public
		postRoutes.GET("", blogHandler.ListPosts)
		postRoutes.GET("/:id", blogHandler.GetPostByID)

		// Protected routes 
		postRoutes.Use(services.AuthMiddleware(jwtService)) // updated name
		{
			postRoutes.POST("", blogHandler.CreatePost)
			postRoutes.PUT("/:id", blogHandler.UpdatePost)
			postRoutes.DELETE("/:id", blogHandler.DeletePost)

			postRoutes.POST("/:id/like", blogHandler.LikePost)
			postRoutes.POST("/:id/dislike", blogHandler.DislikePost)
		}
	}

	// Admin routes
	adminGroup := router.Group("/admin")
	adminGroup.Use(services.AuthMiddleware(jwtService))
	adminGroup.Use(services.RoleAuthorization("admin"))
	{
		adminGroup.PUT("/users/:id/promote", userManagementHandler.PromoteUser)
		adminGroup.PUT("/users/:id/demote", userManagementHandler.DemoteUser)
		adminGroup.GET("/users/:id", userManagementHandler.GetUserByID)
	}

	return router
}
