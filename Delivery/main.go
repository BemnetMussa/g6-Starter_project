package main

import (
	"context"
	"log"
	"os"
	"starter_project/Delivery/handlers"
	"starter_project/Delivery/routers"
	"starter_project/Infrastructure/mongodb/repositories"
	"starter_project/Usecases"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "blogdb_dev" //  database name
	}

	// --- DATABASE CONNECTION ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	log.Println("Connected to MongoDB!")

	// jwtService := services.NewJWTService() //service part

	// Repository Layer
	// userRepo := repositories.NewUserRepository(db) //user_part
	blogRepo := repositories.NewBlogRepository(db)
	interactionRepo := repositories.NewBlogInteractionRepository(db)

	// Usecase Layer
	// userUsecase := Usecases.NewUserUsecase(userRepo, nil, nil, nil) // user_part
	blogUsecase := Usecases.NewBlogUsecase(blogRepo, interactionRepo, userRepo)

	// Delivery Layer (Handlers)
	var userHandler *handlers.UserHandler
	blogHandler := handlers.NewBlogHandler(blogUsecase)

	// --- SETUP ROUTER AND START SERVER ---
	router := routers.SetupRouter(userHandler, blogHandler, jwtService)

	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
