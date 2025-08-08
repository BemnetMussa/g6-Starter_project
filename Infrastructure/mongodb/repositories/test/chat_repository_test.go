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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatTestSuite struct {
	client         *mongo.Client
	database       *mongo.Database
	chatCollection *mongo.Collection
	chatRepo       entities.ChatRepository
	config         *TestConfig
}

func setupChatTestSuite(t *testing.T) *ChatTestSuite {
	config := GetTestConfig()
	client, database, _ := SetupTestDatabase(t, config)
	
	// Create collection for chat testing
	chatCollection := database.Collection("chats")

	// Clear collection before each test
	_, err := chatCollection.DeleteMany(context.TODO(), bson.M{})
	require.NoError(t, err)

	// Create repository
	chatRepo := repositories.NewChatRepository(chatCollection)

	return &ChatTestSuite{
		client:         client,
		database:       database,
		chatCollection: chatCollection,
		chatRepo:       chatRepo,
		config:         config,
	}
}

func (ts *ChatTestSuite) teardown(t *testing.T) {
	CleanupTestDatabase(t, ts.client, ts.database)
}

func createTestChat(userID primitive.ObjectID) *entities.Chat {
	now := time.Now()
	return &entities.Chat{
		UserID:    userID,
		Request:   "Hello, how are you?",
		Response:  "I'm doing well, thank you for asking!",
		Tokens:    50,
		CreatedAt: now,
	}
}

func createTestChatWithCustomFields(userID primitive.ObjectID, request, response string, tokens int) *entities.Chat {
	now := time.Now()
	return &entities.Chat{
		UserID:    userID,
		Request:   request,
		Response:  response,
		Tokens:    tokens,
		CreatedAt: now,
	}
}

func TestChatRepository_CreateChat(t *testing.T) {
	ts := setupChatTestSuite(t)
	defer ts.teardown(t)

	t.Run("should create chat successfully", func(t *testing.T) {
		userID := primitive.NewObjectID()
		chat := createTestChat(userID)

		createdChat, err := ts.chatRepo.CreateChat(chat)

		assert.NoError(t, err)
		assert.NotNil(t, createdChat)
		assert.NotEmpty(t, createdChat.ID)
		assert.Equal(t, userID, createdChat.UserID)
		assert.Equal(t, chat.Request, createdChat.Request)
		assert.Equal(t, chat.Response, createdChat.Response)
		assert.Equal(t, chat.Tokens, createdChat.Tokens)
	})

	t.Run("should create chat with custom fields", func(t *testing.T) {
		userID := primitive.NewObjectID()
		request := "What is the capital of France?"
		response := "The capital of France is Paris."
		tokens := 25
		
		chat := createTestChatWithCustomFields(userID, request, response, tokens)

		createdChat, err := ts.chatRepo.CreateChat(chat)

		assert.NoError(t, err)
		assert.NotNil(t, createdChat)
		assert.Equal(t, userID, createdChat.UserID)
		assert.Equal(t, request, createdChat.Request)
		assert.Equal(t, response, createdChat.Response)
		assert.Equal(t, tokens, createdChat.Tokens)
	})
}

func TestChatRepository_GetChatByID(t *testing.T) {
	ts := setupChatTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find chat by ID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		chat := createTestChat(userID)
		
		createdChat, err := ts.chatRepo.CreateChat(chat)
		require.NoError(t, err)

		foundChat, err := ts.chatRepo.GetChatByID(createdChat.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, foundChat)
		assert.Equal(t, createdChat.ID, foundChat.ID)
		assert.Equal(t, createdChat.UserID, foundChat.UserID)
		assert.Equal(t, createdChat.Request, foundChat.Request)
		assert.Equal(t, createdChat.Response, foundChat.Response)
	})

	t.Run("should return error for invalid ID format", func(t *testing.T) {
		_, err := ts.chatRepo.GetChatByID("invalid-id")

		assert.Error(t, err)
		assert.Equal(t, "invalid chat ID", err.Error())
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		validID := primitive.NewObjectID().Hex()
		_, err := ts.chatRepo.GetChatByID(validID)

		assert.Error(t, err)
		assert.Equal(t, "chat not found", err.Error())
	})
}

func TestChatRepository_GetChatsByUserID(t *testing.T) {
	ts := setupChatTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find chats by user ID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		
		// Create multiple chats for the same user
		chat1 := createTestChat(userID)
		chat2 := createTestChatWithCustomFields(userID, "Second request", "Second response", 30)
		chat3 := createTestChatWithCustomFields(userID, "Third request", "Third response", 40)

		_, err := ts.chatRepo.CreateChat(chat1)
		require.NoError(t, err)
		_, err = ts.chatRepo.CreateChat(chat2)
		require.NoError(t, err)
		_, err = ts.chatRepo.CreateChat(chat3)
		require.NoError(t, err)

		foundChats, err := ts.chatRepo.GetChatsByUserID(userID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, foundChats)
		assert.Len(t, foundChats, 3)
		
		// Verify all chats belong to the same user
		for _, chat := range foundChats {
			assert.Equal(t, userID, chat.UserID)
		}
	})

	t.Run("should return empty slice for user with no chats", func(t *testing.T) {
		userID := primitive.NewObjectID()

		foundChats, err := ts.chatRepo.GetChatsByUserID(userID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, foundChats)
		assert.Len(t, foundChats, 0)
	})

	t.Run("should return error for invalid user ID format", func(t *testing.T) {
		_, err := ts.chatRepo.GetChatsByUserID("invalid-user-id")

		assert.Error(t, err)
		assert.Equal(t, "invalid user ID", err.Error())
	})

	t.Run("should not return chats from other users", func(t *testing.T) {
		user1ID := primitive.NewObjectID()
		user2ID := primitive.NewObjectID()
		
		// Create chat for user1
		chat1 := createTestChat(user1ID)
		_, err := ts.chatRepo.CreateChat(chat1)
		require.NoError(t, err)

		// Create chat for user2
		chat2 := createTestChat(user2ID)
		_, err = ts.chatRepo.CreateChat(chat2)
		require.NoError(t, err)

		// Get chats for user1 only
		foundChats, err := ts.chatRepo.GetChatsByUserID(user1ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, foundChats)
		assert.Len(t, foundChats, 1)
		assert.Equal(t, user1ID, foundChats[0].UserID)
	})
}

func TestChatRepository_DeleteChat(t *testing.T) {
	ts := setupChatTestSuite(t)
	defer ts.teardown(t)

	t.Run("should delete chat successfully", func(t *testing.T) {
		userID := primitive.NewObjectID()
		chat := createTestChat(userID)
		
		createdChat, err := ts.chatRepo.CreateChat(chat)
		require.NoError(t, err)

		err = ts.chatRepo.DeleteChat(createdChat.ID.Hex())

		assert.NoError(t, err)

		// Verify chat is deleted
		_, err = ts.chatRepo.GetChatByID(createdChat.ID.Hex())
		assert.Error(t, err)
		assert.Equal(t, "chat not found", err.Error())
	})

	t.Run("should return error for invalid chat ID format", func(t *testing.T) {
		err := ts.chatRepo.DeleteChat("invalid-id")

		assert.Error(t, err)
		assert.Equal(t, "invalid chat ID", err.Error())
	})

	t.Run("should handle deletion of non-existent chat", func(t *testing.T) {
		validID := primitive.NewObjectID().Hex()
		err := ts.chatRepo.DeleteChat(validID)

		assert.NoError(t, err) // DeleteOne doesn't return error if no document found
	})
}

func TestChatRepository_Integration(t *testing.T) {
	ts := setupChatTestSuite(t)
	defer ts.teardown(t)

	t.Run("should handle complete chat lifecycle", func(t *testing.T) {
		userID := primitive.NewObjectID()
		chat := createTestChat(userID)

		// Create chat
		createdChat, err := ts.chatRepo.CreateChat(chat)
		assert.NoError(t, err)
		assert.NotNil(t, createdChat)

		// Find chat by ID
		foundChat, err := ts.chatRepo.GetChatByID(createdChat.ID.Hex())
		assert.NoError(t, err)
		assert.NotNil(t, foundChat)
		assert.Equal(t, createdChat.ID, foundChat.ID)

		// Find chats by user ID
		foundChats, err := ts.chatRepo.GetChatsByUserID(userID.Hex())
		assert.NoError(t, err)
		assert.Len(t, foundChats, 1)
		assert.Equal(t, createdChat.ID, foundChats[0].ID)

		// Delete chat
		err = ts.chatRepo.DeleteChat(createdChat.ID.Hex())
		assert.NoError(t, err)

		// Verify chat is deleted
		_, err = ts.chatRepo.GetChatByID(createdChat.ID.Hex())
		assert.Error(t, err)
		assert.Equal(t, "chat not found", err.Error())

		// Verify user has no chats
		foundChats, err = ts.chatRepo.GetChatsByUserID(userID.Hex())
		assert.NoError(t, err)
		assert.Len(t, foundChats, 0)
	})

	t.Run("should handle multiple users with multiple chats", func(t *testing.T) {
		user1ID := primitive.NewObjectID()
		user2ID := primitive.NewObjectID()

		// Create chats for user1
		chat1 := createTestChat(user1ID)
		chat2 := createTestChatWithCustomFields(user1ID, "User1 second request", "User1 second response", 35)
		
		_, err := ts.chatRepo.CreateChat(chat1)
		require.NoError(t, err)
		_, err = ts.chatRepo.CreateChat(chat2)
		require.NoError(t, err)

		// Create chat for user2
		chat3 := createTestChat(user2ID)
		_, err = ts.chatRepo.CreateChat(chat3)
		require.NoError(t, err)

		// Get chats for user1
		user1Chats, err := ts.chatRepo.GetChatsByUserID(user1ID.Hex())
		assert.NoError(t, err)
		assert.Len(t, user1Chats, 2)

		// Get chats for user2
		user2Chats, err := ts.chatRepo.GetChatsByUserID(user2ID.Hex())
		assert.NoError(t, err)
		assert.Len(t, user2Chats, 1)

		// Delete one chat from user1
		err = ts.chatRepo.DeleteChat(user1Chats[0].ID.Hex())
		assert.NoError(t, err)

		// Verify user1 now has 1 chat
		user1Chats, err = ts.chatRepo.GetChatsByUserID(user1ID.Hex())
		assert.NoError(t, err)
		assert.Len(t, user1Chats, 1)

		// Verify user2 still has 1 chat
		user2Chats, err = ts.chatRepo.GetChatsByUserID(user2ID.Hex())
		assert.NoError(t, err)
		assert.Len(t, user2Chats, 1)
	})
} 