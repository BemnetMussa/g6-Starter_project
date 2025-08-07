package routers

import (
	"g6_starter_project/Delivery/handlers"
	"g6_starter_project/Infrastructure/services"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userUsecase *usecases.UserUsecase,
	passwordResetUsecase *usecases.PasswordResetUsecase,
	userManagementUsecase *usecases.UserManagementUsecase,
	blogHandler *handlers.BlogHandler,
	authSvc services.JWTServiceInterface,

	userProfileHandler *handlers.UserProfileHandler,
	commentHandler *handlers.CommentHandler,

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

	logoutRoutes := router.Group("")
	logoutRoutes.Use(services.AuthMiddleware(authSvc))
	{
		logoutRoutes.POST("/logout", userHandler.Logout)
	}

	// Profile routes (authentication required)
	profileRoutes := router.Group("/profile")
	profileRoutes.Use(services.GinAuthMiddleware(authSvc))
	{
		profileRoutes.GET("/me", userProfileHandler.GetMyProfile)
		profileRoutes.PUT("/me", userProfileHandler.UpdateMyProfile)
	}

	// Blog routes
	postRoutes := router.Group("/blog")
	{
		// Public
		postRoutes.GET("", blogHandler.ListPosts)
		postRoutes.GET("/:id", blogHandler.GetPostByID)

		// Protected routes
		protectedPostRoutes := postRoutes.Group("")
		protectedPostRoutes.Use(services.GinAuthMiddleware(authSvc))

		{
			protectedPostRoutes.POST("", blogHandler.CreatePost)
			protectedPostRoutes.PUT("/:id", blogHandler.UpdatePost)
			protectedPostRoutes.DELETE("/:id", blogHandler.DeletePost)

			protectedPostRoutes.POST("/:id/like", blogHandler.LikePost)
			protectedPostRoutes.POST("/:id/dislike", blogHandler.DislikePost)
			protectedPostRoutes.POST("/:id/comments", commentHandler.CreateComment)
		}
	}

	// Admin routes
	adminGroup := router.Group("/admin")
	adminGroup.Use(services.GinAuthMiddleware(authSvc))
	adminGroup.Use(services.GinRoleAuthorization("admin"))
	{
		adminGroup.PUT("/users/:id/promote", userManagementHandler.PromoteUser)
		adminGroup.PUT("/users/:id/demote", userManagementHandler.DemoteUser)
		adminGroup.GET("/users/:id", userManagementHandler.GetUserByID)
	}

	return router
}
