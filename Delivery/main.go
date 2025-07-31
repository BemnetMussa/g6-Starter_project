package main

import (
	"context"
	"log"
	"os"

	"g6_starter_project/Delivery/routers"
	"g6_starter_project/Infrastructure/db"
	"g6_starter_project/Infrastructure/mongodb/repositories"
	"g6_starter_project/Infrastructure/services"
	usecases "g6_starter_project/Usecases"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	// Get database name from environment
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "blog_api" // default fallback
	}

	// Get port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // default fallback
	}

	// Initialize MongoDB client
	client, err := db.ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Get database
	database := client.Database(dbName)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.Collection("users"))
	tokenRepo := repositories.NewTokenRepository(database.Collection("token"))

	// Initialize services
	jwtService := services.NewJWTService(os.Getenv("JWT_SECRET")) 

	// Initialize usecases
	tokenUsecase := usecases.NewTokenUsecase(tokenRepo, jwtService)
	userUsecase := usecases.NewUserUsecase(userRepo, tokenUsecase)

	// Setup router
	router := routers.SetupRouter(userUsecase)	

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
