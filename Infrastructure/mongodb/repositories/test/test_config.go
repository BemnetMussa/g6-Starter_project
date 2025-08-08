package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"g6_starter_project/Domain/entities"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestConfig holds test configuration
type TestConfig struct {
	MongoURI      string
	DatabaseName  string
	CollectionName string
}

// GetTestConfig returns test configuration from environment or defaults
func GetTestConfig() *TestConfig {
	// Try multiple paths for .env file
	envPaths := []string{
		".env",
		"../.env", 
		"../../.env",
		"../../../.env",
		"../../../../.env",
	}
	
	var err error
	for _, path := range envPaths {
		err = godotenv.Load(path)
		if err == nil {
			fmt.Printf("Loaded .env from: %s\n", path)
			break
		}
	}
	
	if err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Debug: Check if environment variable is set
	mongoURI := os.Getenv("MONGODB_URI")
	fmt.Printf("Environment MONGODB_URI: %s\n", mongoURI)

	config := &TestConfig{
		MongoURI:       getEnvOrDefault("MONGODB_URI", "mongodb+srv://leulgedion:kT5JsmzjYL8hVBrc@cluster0.1y2cmpf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"),
		DatabaseName:   getEnvOrDefault("MONGODB_DATABASE", "test_blog_api"),
		CollectionName: getEnvOrDefault("MONGODB_COLLECTION", "users"),
	}

	fmt.Printf("Test Config - MongoDB URI: %s\n", config.MongoURI)
	fmt.Printf("Test Config - Database: %s\n", config.DatabaseName)
	fmt.Printf("Test Config - Collection: %s\n", config.CollectionName)

	return config
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTestDatabase creates a test database connection
func SetupTestDatabase(t *testing.T, config *TestConfig) (*mongo.Client, *mongo.Database, *mongo.Collection) {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	require.NoError(t, err)

	// Ping the database
	err = client.Ping(context.TODO(), nil)
	require.NoError(t, err)

	// Create test database and collection
	database := client.Database(config.DatabaseName)
	collection := database.Collection(config.CollectionName)

	// Clear the collection before each test
	_, err = collection.DeleteMany(context.TODO(), bson.M{})
	require.NoError(t, err)

	return client, database, collection
}

// CleanupTestDatabase cleans up test database
func CleanupTestDatabase(t *testing.T, client *mongo.Client, database *mongo.Database) {
	if client != nil {
		// Don't drop the database - just disconnect
		// The user doesn't have permission to drop databases
		err := client.Disconnect(context.TODO())
		require.NoError(t, err)
	}
}

// CreateTestUser creates a test user with default values
func CreateTestUser() *entities.User {
	now := time.Now()
	return &entities.User{
		FullName:   "Test User",
		Username:   "testuser",
		Email:      "test@example.com",
		Password:   "hashedpassword",
		Role:       "user",
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// CreateTestUserWithCustomFields creates a test user with custom fields
func CreateTestUserWithCustomFields(fullName, username, email string) *entities.User {
	user := CreateTestUser()
	user.FullName = fullName
	user.Username = username
	user.Email = email
	return user
}

// CreateVerifiedUser creates a verified test user
func CreateVerifiedUser() *entities.User {
	user := CreateTestUser()
	user.IsVerified = true
	return user
}

// CreateAdminUser creates an admin test user
func CreateAdminUser() *entities.User {
	user := CreateTestUser()
	user.Role = "admin"
	user.IsVerified = true
	return user
}

// CreateUserWithContactInfo creates a test user with contact information
func CreateUserWithContactInfo() *entities.User {
	user := CreateTestUser()
	phone := "+1234567890"
	address := "123 Test Street, Test City"
	user.ContactInfo = &entities.ContactInfo{
		Phone:   &phone,
		Address: &address,
	}
	return user
}

// CreateUserWithProfileImage creates a test user with profile image
func CreateUserWithProfileImage() *entities.User {
	user := CreateTestUser()
	imageURL := "https://example.com/profile.jpg"
	user.ProfileImage = &imageURL
	return user
}

// CreateUserWithBio creates a test user with bio
func CreateUserWithBio() *entities.User {
	user := CreateTestUser()
	bio := "This is a test bio for the user"
	user.Bio = &bio
	return user
}

// CreateUserWithResetToken creates a test user with reset token
func CreateUserWithResetToken() *entities.User {
	user := CreateTestUser()
	token := "test-reset-token-123"
	expiresAt := time.Now().Add(15 * time.Minute)
	user.ResetToken = &token
	user.ResetTokenExpiresAt = &expiresAt
	return user
}

// Helper function to create string pointer
func StringPtr(s string) *string {
	return &s
}

// Helper function to create time pointer
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Helper function to create ObjectID pointer
func ObjectIDPtr(id primitive.ObjectID) *primitive.ObjectID {
	return &id
}

// AssertUserFields asserts common user fields
func AssertUserFields(t *testing.T, expected, actual *entities.User) {
	require.NotNil(t, actual)
	require.NotEmpty(t, actual.ID)
	require.Equal(t, expected.FullName, actual.FullName)
	require.Equal(t, expected.Username, actual.Username)
	require.Equal(t, expected.Email, actual.Email)
	require.Equal(t, expected.Role, actual.Role)
	require.Equal(t, expected.IsVerified, actual.IsVerified)
}

// AssertUserNotExists asserts that a user doesn't exist
func AssertUserNotExists(t *testing.T, repo entities.UserRepository, userID string) {
	_, err := repo.GetUserByID(userID)
	require.Error(t, err)
	require.Equal(t, "user not found", err.Error())
}

// AssertUserExists asserts that a user exists
func AssertUserExists(t *testing.T, repo entities.UserRepository, userID string) *entities.User {
	user, err := repo.GetUserByID(userID)
	require.NoError(t, err)
	require.NotNil(t, user)
	return user
} 