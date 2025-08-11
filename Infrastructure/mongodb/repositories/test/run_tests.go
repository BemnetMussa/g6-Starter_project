package test

import (
	"fmt"
	"os"
	"testing"
)

// TestMain runs before all tests in this package
func TestMain(m *testing.M) {
	// Setup any global test configuration here
	fmt.Println("Starting repository tests...")
	
	// Run the tests
	exitCode := m.Run()
	
	// Cleanup after all tests
	fmt.Println("Repository tests completed.")
	
	os.Exit(exitCode)
}

// RunAllRepositoryTests is a helper function to run all repository tests
func RunAllRepositoryTests(t *testing.T) {
	t.Run("UserRepository", func(t *testing.T) {
		// User repository tests are already defined in user_repository_test.go
		// This is just a wrapper to group them
	})
	
	t.Run("BlogRepository", func(t *testing.T) {
		// Blog repository tests are already defined in blog_repository_test.go
		// This is just a wrapper to group them
	})
	
	t.Run("TokenRepository", func(t *testing.T) {
		// Token repository tests are already defined in token_repository_test.go
		// This is just a wrapper to group them
	})
	
	t.Run("ChatRepository", func(t *testing.T) {
		// Chat repository tests are already defined in chat_repository_test.go
		// This is just a wrapper to group them
	})
} 