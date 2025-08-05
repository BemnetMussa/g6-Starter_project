package routers

import (
	"g6_starter_project/Delivery/handlers"
	"g6_starter_project/Infrastructure/services"
	usecases "g6_starter_project/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userUsecase *usecases.UserUsecase,
	blogHandler *handlers.BlogHandler,
	authSvc services.JWTServiceInterface,
) *gin.Engine {

	router := gin.Default()
	userHandler := handlers.NewUserHandler(userUsecase)

	// Public routes
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	// Blog routes
	postRoutes := router.Group("/blog")
	{
		// Public
		postRoutes.GET("", blogHandler.ListPosts)
		postRoutes.GET("/:id", blogHandler.GetPostByID)

		// Protected
		protectedPostRoutes := postRoutes.Group("")
		protectedPostRoutes.Use(services.AuthMiddleware(authSvc))
		{
			protectedPostRoutes.POST("", blogHandler.CreatePost)
			protectedPostRoutes.PUT("/:id", blogHandler.UpdatePost)
			protectedPostRoutes.DELETE("/:id", blogHandler.DeletePost)

			protectedPostRoutes.POST("/:id/like", blogHandler.LikePost)
			protectedPostRoutes.POST("/:id/dislike", blogHandler.DislikePost)
		}
	}

	return router
}
