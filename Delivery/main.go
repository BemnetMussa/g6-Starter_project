package main

import (
	"context"
	"log"
	"os"

	"g6_starter_project/Delivery/routers"
	"g6_starter_project/Infrastructure/db"
	"g6_starter_project/Infrastructure/mongodb/repositories"
	"g6_starter_project/Delivery/handlers"
	"g6_starter_project/Infrastructure/services"
	 Usecases "g6_starter_project/Usecases"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	LoadEnvVariables()

	mongoURI, databaseName, serverPort := GetAppConfig()

	mongoClient := ConnectToMongoDB(mongoURI)
	defer mongoClient.Disconnect(context.TODO())

	database := mongoClient.Database(databaseName)

	// Services
	jwtService := services.NewJWTService(os.Getenv("JWT_SECRET"))
	emailService := services.NewEmailService()
	rateLimiter := services.NewRateLimiter()
	rateLimiter.StartCleanup()

	// Repositories
	userRepository := repositories.NewUserRepository(database.Collection("users"))
	tokenRepository := repositories.NewTokenRepository(database.Collection("token"))
	blogRepository := repositories.NewBlogRepository(database)
	interactionRepository := repositories.NewBlogInteractionRepository(database)

	// UseCases
	tokenUseCase := Usecases.NewTokenUsecase(tokenRepository, jwtService)
	userUseCase := Usecases.NewUserUsecase(userRepository, tokenUseCase)
	passwordResetUseCase := Usecases.NewPasswordResetUsecase(userRepository, jwtService, emailService, rateLimiter)
	userManagementUseCase := Usecases.NewUserManagementUsecase(userRepository)
	blogUseCase := Usecases.NewBlogUsecase(blogRepository, interactionRepository, userRepository)

	// Handlers
	blogHandler := handlers.NewBlogHandler(blogUseCase)

	// Router
	router := routers.SetupRouter(userUseCase, passwordResetUseCase, userManagementUseCase, blogHandler, jwtService)

	log.Printf("Server running on port %s", serverPort)
	if err := router.Run(":" + serverPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}





func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
}

func GetAppConfig() (mongoURI string, dbName string, port string) {
	mongoURI = os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("Environment variable MONGODB_URI is required")
	}

	dbName = os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "blog_api"
	}

	port = os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return
}

func ConnectToMongoDB(uri string) *mongo.Client {
	client, err := db.ConnectMongoDB(uri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	return client
}
