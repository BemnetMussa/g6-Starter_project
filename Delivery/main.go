package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"g6-Starter_project/Delivery/routers"
	"g6-Starter_project/Infrastructure/db"
	"g6-Starter_project/Infrastructure/mongodb/repositories"
	"g6-Starter_project/Usecases"
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

	// Initialize usecases
	userUsecase := usecases.NewUserUsecase(userRepo)

	// Setup router
	router := routers.SetupRouter(userUsecase)

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
