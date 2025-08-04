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
	
	// Public routes (no authentication required)
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.POST("/forgot-password", userHandler.ForgotPassword)
	router.POST("/reset-password", userHandler.ResetPassword)
	
	// Blog routes
	postRoutes := router.Group("/blog")
	{
		// Public blog routes
		postRoutes.GET("", blogHandler.ListPosts)
		postRoutes.GET("/:id", blogHandler.GetPostByID)

		// Protected blog routes
		protectedPostRoutes := postRoutes.Group("")
		protectedPostRoutes.Use(services.GinAuthMiddleware(jwtService)) 
		{
			protectedPostRoutes.POST("", blogHandler.CreatePost)
			protectedPostRoutes.PUT("/:id", blogHandler.UpdatePost)
			protectedPostRoutes.DELETE("/:id", blogHandler.DeletePost)

			protectedPostRoutes.POST("/:id/like", blogHandler.LikePost)
			protectedPostRoutes.POST("/:id/dislike", blogHandler.DislikePost)
		}
	}
	
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
