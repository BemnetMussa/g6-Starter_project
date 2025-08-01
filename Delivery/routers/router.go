package routers

import (
	"g6_starter_project/Infrastructure/services"
	"g6_starter_project/Delivery/handlers"
	usecases "g6_starter_project/Usecases"


	"github.com/gin-gonic/gin"
)

func SetupRouter(
  userUsecase *usecases.UserUsecase
	blogHandler *handlers.BlogHandler,
	// user handler
	// jwtService services.IJWTService,
) *gin.Engine {

	r := gin.Default()
  userHandler := handlers.NewUserHandler(userUsecase)


  router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	postRoutes := r.Group("/blog")
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
