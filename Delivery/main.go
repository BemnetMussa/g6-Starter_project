package main

import (
	"context"
	"log"
	"g6-Starter_project/Delivery/routers"
	"g6-Starter_project/Infrastructure/db"
	"g6-Starter_project/Infrastructure/mongodb/repositories"
	"g6-Starter_project/Usecases"
)

func main() {
	// Initialize MongoDB client
	client, err := db.ConnectMongoDB("mongodb+srv://leulgedion:kT5JsmzjYL8hVBrc@cluster0.1y2cmpf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Get database
	database := client.Database("blog_api")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.Collection("users"))

	// Initialize usecases
	userUsecase := usecases.NewUserUsecase(userRepo)

	// Setup router
	router := routers.SetupRouter(userUsecase)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
