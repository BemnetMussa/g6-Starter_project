package routers

import (
	"starter_project/Delivery/handlers"
	"starter_project/Infrastructure/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	blogHandler *handlers.BlogHandler,
	// user handler
	// jwtService services.IJWTService,
) *gin.Engine {

	r := gin.Default()
	// User part
	// authRoutes := r.Group("/auth")
	// {
	// 	authRoutes.POST("/register", userHandler.Register)
	// 	authRoutes.POST("/login", userHandler.Login)
	// }

	postRoutes := r.Group("/posts")
	{
		// These two routes are PUBLIC
		postRoutes.GET("", blogHandler.ListPosts)
		postRoutes.GET("/:id", blogHandler.GetPostByID)

		// These routes are PROTECTED
		protectedPostRoutes := postRoutes.Use(services.AuthMiddleware(jwtService))
		{
			protectedPostRoutes.POST("", blogHandler.CreatePost)
			protectedPostRoutes.PUT("/:id", blogHandler.UpdatePost)
			protectedPostRoutes.DELETE("/:id", blogHandler.DeletePost)

			// Routes for popularity features
			protectedPostRoutes.POST("/:id/like", blogHandler.LikePost)
			protectedPostRoutes.POST("/:id/dislike", blogHandler.DislikePost)
		}
	}

	return r
}
