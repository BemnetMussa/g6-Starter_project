package test

import (
	"context"
	"testing"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/mongodb/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenTestSuite struct {
	client         *mongo.Client
	database       *mongo.Database
	tokenCollection *mongo.Collection
	tokenRepo      *repositories.TokenRepository
	config         *TestConfig
}

func setupTokenTestSuite(t *testing.T) *TokenTestSuite {
	config := GetTestConfig()
	client, database, _ := SetupTestDatabase(t, config)
	
	// Create collection for token testing
	tokenCollection := database.Collection("tokens")

	// Clear collection before each test
	_, err := tokenCollection.DeleteMany(context.TODO(), bson.M{})
	require.NoError(t, err)

	// Create repository
	tokenRepo := repositories.NewTokenRepository(tokenCollection)

	return &TokenTestSuite{
		client:         client,
		database:       database,
		tokenCollection: tokenCollection,
		tokenRepo:      tokenRepo,
		config:         config,
	}
}

func (ts *TokenTestSuite) teardown(t *testing.T) {
	CleanupTestDatabase(t, ts.client, ts.database)
}

func createTestToken(userID string) *entities.Token {
	now := time.Now().UTC().Truncate(time.Second)
	return &entities.Token{
		UserID:       userID,
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    now.Add(24 * time.Hour),
		CreatedAt:    now,
	}
}

func createTestTokenWithCustomFields(userID, accessToken, refreshToken string, expiresAt time.Time) *entities.Token {
	now := time.Now().UTC().Truncate(time.Second)
	return &entities.Token{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Truncate(time.Second),
		CreatedAt:    now,
	}
}

func TestTokenRepository_Create(t *testing.T) {
	ts := setupTokenTestSuite(t)
	defer ts.teardown(t)

	t.Run("should create token successfully", func(t *testing.T) {
		userID := "test-user-id"
		token := createTestToken(userID)

		err := ts.tokenRepo.Create(context.TODO(), token)

		assert.NoError(t, err)
		assert.NotEmpty(t, token.ID)
	})

	t.Run("should create token with custom fields", func(t *testing.T) {
		userID := "test-user-id-2"
		accessToken := "custom-access-token"
		refreshToken := "custom-refresh-token"
		expiresAt := time.Now().Add(48 * time.Hour)
		
		token := createTestTokenWithCustomFields(userID, accessToken, refreshToken, expiresAt)

		err := ts.tokenRepo.Create(context.TODO(), token)

		assert.NoError(t, err)
		assert.NotEmpty(t, token.ID)
		assert.Equal(t, userID, token.UserID)
		assert.Equal(t, accessToken, token.AccessToken)
		assert.Equal(t, refreshToken, token.RefreshToken)
		assert.WithinDuration(t, expiresAt, token.ExpiresAt, time.Second)
	})
}

func TestTokenRepository_FindByUserID(t *testing.T) {
	ts := setupTokenTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find token by user ID", func(t *testing.T) {
		userID := "test-user-id"
		token := createTestToken(userID)

		err := ts.tokenRepo.Create(context.TODO(), token)
		require.NoError(t, err)

		foundToken, err := ts.tokenRepo.FindByUserID(context.TODO(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, foundToken)
		assert.Equal(t, userID, foundToken.UserID)
		assert.Equal(t, token.AccessToken, foundToken.AccessToken)
		assert.Equal(t, token.RefreshToken, foundToken.RefreshToken)
		// Compare time with tolerance for precision differences
		assert.WithinDuration(t, token.ExpiresAt, foundToken.ExpiresAt, time.Second)
	})

	t.Run("should return error for non-existent user ID", func(t *testing.T) {
		nonExistentUserID := "non-existent-user-id"

		_, err := ts.tokenRepo.FindByUserID(context.TODO(), nonExistentUserID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token not found for user")
	})
}

func TestTokenRepository_Update(t *testing.T) {
	ts := setupTokenTestSuite(t)
	defer ts.teardown(t)

	t.Run("should update token successfully", func(t *testing.T) {
		userID := "test-user-id"
		token := createTestToken(userID)

		err := ts.tokenRepo.Create(context.TODO(), token)
		require.NoError(t, err)

		// Update token
		newAccessToken := "updated-access-token"
		newRefreshToken := "updated-refresh-token"
		newExpiresAt := time.Now().Add(72 * time.Hour)

		update := bson.M{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"expires_at":    newExpiresAt,
		}

		err = ts.tokenRepo.Update(context.TODO(), userID, update)

		assert.NoError(t, err)

		// Verify update
		updatedToken, err := ts.tokenRepo.FindByUserID(context.TODO(), userID)
		assert.NoError(t, err)
		assert.Equal(t, newAccessToken, updatedToken.AccessToken)
		assert.Equal(t, newRefreshToken, updatedToken.RefreshToken)
		assert.WithinDuration(t, newExpiresAt, updatedToken.ExpiresAt, time.Second)
	})

	t.Run("should update only specific fields", func(t *testing.T) {
		userID := "test-user-id-2"
		token := createTestToken(userID)

		err := ts.tokenRepo.Create(context.TODO(), token)
		require.NoError(t, err)

		// Update only access token
		newAccessToken := "partial-updated-access-token"
		update := bson.M{
			"access_token": newAccessToken,
		}

		err = ts.tokenRepo.Update(context.TODO(), userID, update)

		assert.NoError(t, err)

		// Verify only access token was updated
		updatedToken, err := ts.tokenRepo.FindByUserID(context.TODO(), userID)
		assert.NoError(t, err)
		assert.Equal(t, newAccessToken, updatedToken.AccessToken)
		assert.Equal(t, token.RefreshToken, updatedToken.RefreshToken) // Should remain unchanged
		assert.WithinDuration(t, token.ExpiresAt, updatedToken.ExpiresAt, time.Second) // Should remain unchanged
	})

	t.Run("should return error for non-existent user ID", func(t *testing.T) {
		nonExistentUserID := "non-existent-user-id"
		update := bson.M{
			"access_token": "new-token",
		}

		err := ts.tokenRepo.Update(context.TODO(), nonExistentUserID, update)

		assert.NoError(t, err) // UpdateOne doesn't return error if no document found
	})
}

func TestTokenRepository_DeleteByUserID(t *testing.T) {
	ts := setupTokenTestSuite(t)
	defer ts.teardown(t)

	t.Run("should delete token by user ID", func(t *testing.T) {
		userID := "test-user-id"
		token := createTestToken(userID)

		err := ts.tokenRepo.Create(context.TODO(), token)
		require.NoError(t, err)

		err = ts.tokenRepo.DeleteByUserID(context.TODO(), userID)

		assert.NoError(t, err)

		// Verify token is deleted
		_, err = ts.tokenRepo.FindByUserID(context.TODO(), userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token not found for user")
	})

	t.Run("should handle deletion of non-existent user ID", func(t *testing.T) {
		nonExistentUserID := "non-existent-user-id"

		err := ts.tokenRepo.DeleteByUserID(context.TODO(), nonExistentUserID)

		assert.NoError(t, err) // DeleteMany doesn't return error if no documents found
	})

	t.Run("should delete multiple tokens for same user ID", func(t *testing.T) {
		userID := "test-user-id"
		token1 := createTestToken(userID)
		token2 := createTestTokenWithCustomFields(userID, "access-2", "refresh-2", time.Now().Add(24*time.Hour))

		err := ts.tokenRepo.Create(context.TODO(), token1)
		require.NoError(t, err)
		err = ts.tokenRepo.Create(context.TODO(), token2)
		require.NoError(t, err)

		err = ts.tokenRepo.DeleteByUserID(context.TODO(), userID)

		assert.NoError(t, err)

		// Verify all tokens are deleted
		_, err = ts.tokenRepo.FindByUserID(context.TODO(), userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token not found for user")
	})
}

func TestTokenRepository_Integration(t *testing.T) {
	ts := setupTokenTestSuite(t)
	defer ts.teardown(t)

	t.Run("should handle complete token lifecycle", func(t *testing.T) {
		userID := "test-user-id"
		token := createTestToken(userID)

		// Create token
		err := ts.tokenRepo.Create(context.TODO(), token)
		assert.NoError(t, err)

		// Find token
		foundToken, err := ts.tokenRepo.FindByUserID(context.TODO(), userID)
		assert.NoError(t, err)
		assert.NotNil(t, foundToken)

		// Update token
		newAccessToken := "new-access-token"
		update := bson.M{
			"access_token": newAccessToken,
		}
		err = ts.tokenRepo.Update(context.TODO(), userID, update)
		assert.NoError(t, err)

		// Verify update
		updatedToken, err := ts.tokenRepo.FindByUserID(context.TODO(), userID)
		assert.NoError(t, err)
		assert.Equal(t, newAccessToken, updatedToken.AccessToken)

		// Delete token
		err = ts.tokenRepo.DeleteByUserID(context.TODO(), userID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = ts.tokenRepo.FindByUserID(context.TODO(), userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token not found for user")
	})
} 