package test

import (
	"context"
	"testing"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/mongodb/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestSuite struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	repo       entities.UserRepository
	config     *TestConfig
}

func setupTestSuite(t *testing.T) *TestSuite {
	config := GetTestConfig()
	client, database, collection := SetupTestDatabase(t, config)
	
	// Create repository
	repo := repositories.NewUserRepository(collection)

	return &TestSuite{
		client:     client,
		database:   database,
		collection: collection,
		repo:       repo,
		config:     config,
	}
}

func (ts *TestSuite) teardown(t *testing.T) {
	CleanupTestDatabase(t, ts.client, ts.database)
}

func TestCreateUser(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateTestUser()

		createdUser, err := ts.repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		AssertUserFields(t, user, createdUser)
	})

	t.Run("should create verified user", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateVerifiedUser()

		createdUser, err := ts.repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.True(t, createdUser.IsVerified)
	})

	t.Run("should create admin user", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateAdminUser()

		createdUser, err := ts.repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, "admin", createdUser.Role)
		assert.True(t, createdUser.IsVerified)
	})

	t.Run("should create user with contact info", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateUserWithContactInfo()

		createdUser, err := ts.repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.NotNil(t, createdUser.ContactInfo)
		assert.Equal(t, "+1234567890", *createdUser.ContactInfo.Phone)
	})

	t.Run("should create user with profile image", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateUserWithProfileImage()

		createdUser, err := ts.repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.NotNil(t, createdUser.ProfileImage)
		assert.Equal(t, "https://example.com/profile.jpg", *createdUser.ProfileImage)
	})

	t.Run("should create user with bio", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateUserWithBio()

		createdUser, err := ts.repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.NotNil(t, createdUser.Bio)
		assert.Equal(t, "This is a test bio for the user", *createdUser.Bio)
	})

	t.Run("should handle duplicate email", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user1 := CreateTestUserWithCustomFields("User 1", "user1", "test@example.com")
		user2 := CreateTestUserWithCustomFields("User 2", "user2", "test@example.com")

		_, err := ts.repo.CreateUser(user1)
		require.NoError(t, err)

		_, err = ts.repo.CreateUser(user2)

		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
	})

	t.Run("should handle duplicate username", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user1 := CreateTestUserWithCustomFields("User 1", "testuser", "user1@example.com")
		user2 := CreateTestUserWithCustomFields("User 2", "testuser", "user2@example.com")

		_, err := ts.repo.CreateUser(user1)
		require.NoError(t, err)

		_, err = ts.repo.CreateUser(user2)

		assert.Error(t, err)
		assert.Equal(t, "username already exists", err.Error())
	})
}

func TestGetUserByEmail(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find user by email", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		foundUser, err := ts.repo.GetUserByEmail(user.Email)

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("should return error for non-existent email", func(t *testing.T) {
		_, err := ts.repo.GetUserByEmail("nonexistent@example.com")

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("should be case insensitive", func(t *testing.T) {
		user := CreateTestUser()
		_, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		foundUser, err := ts.repo.GetUserByEmail("TEST@EXAMPLE.COM")

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.Email, foundUser.Email)
	})
}

func TestGetUserByUsername(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find user by username", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		foundUser, err := ts.repo.GetUserByUsername(user.Username)

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
	})

	t.Run("should return error for non-existent username", func(t *testing.T) {
		_, err := ts.repo.GetUserByUsername("nonexistent")

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestGetUserByID(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find user by ID", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		foundUser, err := ts.repo.GetUserByID(createdUser.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, createdUser.ID, foundUser.ID)
	})

	t.Run("should return error for invalid ID format", func(t *testing.T) {
		_, err := ts.repo.GetUserByID("invalid-id")

		assert.Error(t, err)
		assert.Equal(t, "invalid user ID", err.Error())
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		validID := primitive.NewObjectID().Hex()
		_, err := ts.repo.GetUserByID(validID)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestGetUserCount(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should return zero for empty collection", func(t *testing.T) {
		count, err := ts.repo.GetUserCount()

		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("should return correct count", func(t *testing.T) {
		// Create multiple users
		user1 := CreateTestUser()
		user2 := CreateTestUserWithCustomFields("Test User 2", "testuser2", "test2@example.com")
		user3 := CreateTestUserWithCustomFields("Test User 3", "testuser3", "test3@example.com")

		_, err := ts.repo.CreateUser(user1)
		require.NoError(t, err)
		_, err = ts.repo.CreateUser(user2)
		require.NoError(t, err)
		_, err = ts.repo.CreateUser(user3)
		require.NoError(t, err)

		count, err := ts.repo.GetUserCount()

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("should update user successfully", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		// Update user
		createdUser.FullName = "Updated Name"
		createdUser.Bio = StringPtr("Updated bio")
		updatedUser, err := ts.repo.UpdateUser(createdUser)

		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, "Updated Name", updatedUser.FullName)
		assert.Equal(t, "Updated bio", *updatedUser.Bio)
		assert.True(t, updatedUser.UpdatedAt.After(createdUser.CreatedAt))
	})

	t.Run("should update user role", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		createdUser.Role = "admin"
		updatedUser, err := ts.repo.UpdateUser(createdUser)

		assert.NoError(t, err)
		assert.Equal(t, "admin", updatedUser.Role)
	})

	t.Run("should update verification status", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		createdUser.IsVerified = true
		updatedUser, err := ts.repo.UpdateUser(createdUser)

		assert.NoError(t, err)
		assert.True(t, updatedUser.IsVerified)
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		ts := setupTestSuite(t)
		defer ts.teardown(t)

		user := CreateTestUser()
		user.ID = primitive.NewObjectID()

		_, err := ts.repo.UpdateUser(user)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestDeleteUser(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should delete user successfully", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		err = ts.repo.DeleteUser(createdUser.ID.Hex())

		assert.NoError(t, err)

		// Verify user is deleted
		AssertUserNotExists(t, ts.repo, createdUser.ID.Hex())
	})

	t.Run("should return error for invalid ID format", func(t *testing.T) {
		err := ts.repo.DeleteUser("invalid-id")

		assert.Error(t, err)
		assert.Equal(t, "invalid user ID", err.Error())
	})
}

func TestUpdateResetToken(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should update reset token successfully", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		resetToken := "test-reset-token"
		expiresAt := time.Now().Add(15 * time.Minute)

		err = ts.repo.UpdateResetToken(createdUser.ID.Hex(), &resetToken, &expiresAt)

		assert.NoError(t, err)

		// Verify token was updated
		foundUser, err := ts.repo.GetUserByResetToken(resetToken)
		assert.NoError(t, err)
		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, resetToken, *foundUser.ResetToken)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		resetToken := "test-reset-token"
		expiresAt := time.Now().Add(15 * time.Minute)

		err := ts.repo.UpdateResetToken("invalid-id", &resetToken, &expiresAt)

		assert.Error(t, err)
		assert.Equal(t, "invalid user ID", err.Error())
	})
}

func TestGetUserByResetToken(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find user by reset token", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		resetToken := "test-reset-token"
		expiresAt := time.Now().Add(15 * time.Minute)

		err = ts.repo.UpdateResetToken(createdUser.ID.Hex(), &resetToken, &expiresAt)
		require.NoError(t, err)

		foundUser, err := ts.repo.GetUserByResetToken(resetToken)

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, resetToken, *foundUser.ResetToken)
	})

	t.Run("should return error for invalid reset token", func(t *testing.T) {
		_, err := ts.repo.GetUserByResetToken("invalid-token")

		assert.Error(t, err)
		assert.Equal(t, "invalid reset token", err.Error())
	})
}

func TestUpdateVerificationStatus(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should update verification status successfully", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		err = ts.repo.UpdateVerificationStatus(createdUser.ID.Hex(), true)

		assert.NoError(t, err)

		// Verify status was updated
		foundUser := AssertUserExists(t, ts.repo, createdUser.ID.Hex())
		assert.True(t, foundUser.IsVerified)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		err := ts.repo.UpdateVerificationStatus("invalid-id", true)

		assert.Error(t, err)
		assert.Equal(t, "invalid user ID", err.Error())
	})
}

func TestFindByName(t *testing.T) {
	ts := setupTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find user by name", func(t *testing.T) {
		user := CreateTestUser()
		createdUser, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		foundUser, err := ts.repo.FindByName(context.TODO(), "Test")

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, createdUser.ID, foundUser.ID)
	})

	t.Run("should return nil for non-existent name", func(t *testing.T) {
		foundUser, err := ts.repo.FindByName(context.TODO(), "NonExistent")

		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("should be case insensitive", func(t *testing.T) {
		user := CreateTestUser()
		_, err := ts.repo.CreateUser(user)
		require.NoError(t, err)

		foundUser, err := ts.repo.FindByName(context.TODO(), "test")

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
	})
} 