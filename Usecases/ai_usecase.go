package usecases

import (
	"errors"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIUsecase interface {
	GenerateBlogContent(userID string, topic string) (*entities.Chat, error)
	SuggestTopics(userID string, category string) (*entities.Chat, error)
	EnhanceContent(userID string, content string) (*entities.Chat, error)
	GetChatHistory(userID string) ([]entities.Chat, error)
	DeleteChat(userID string, chatID string) error
}

type aiUsecase struct {
	aiService      *services.AIService
	chatRepository entities.ChatRepository
	userRepository entities.UserRepository
}

func NewAIUsecase(aiService *services.AIService, chatRepository entities.ChatRepository, userRepository entities.UserRepository) AIUsecase {
	return &aiUsecase{
		aiService:      aiService,
		chatRepository: chatRepository,
		userRepository: userRepository,
	}
}

func (uc *aiUsecase) GenerateBlogContent(userID string, topic string) (*entities.Chat, error) {
	// Validate user exists
	_, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate content using AI service
	content, tokens, err := uc.aiService.GenerateBlogContent(topic)
	if err != nil {
		return nil, err
	}

	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Create chat entity
	chat := &entities.Chat{
		UserID:    userObjectID,
		Request:   topic,
		Response:  content,
		Tokens:    tokens,
		CreatedAt: time.Now(),
	}

	// Save to database
	savedChat, err := uc.chatRepository.CreateChat(chat)
	if err != nil {
		return nil, err
	}

	return savedChat, nil
}

func (uc *aiUsecase) SuggestTopics(userID string, category string) (*entities.Chat, error) {
	// Validate user exists
	_, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate topics using AI service
	topics, tokens, err := uc.aiService.SuggestTopics(category)
	if err != nil {
		return nil, err
	}

	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Create chat entity
	chat := &entities.Chat{
		UserID:    userObjectID,
		Request:   category,
		Response:  topics,
		Tokens:    tokens,
		CreatedAt: time.Now(),
	}

	// Save to database
	savedChat, err := uc.chatRepository.CreateChat(chat)
	if err != nil {
		return nil, err
	}

	return savedChat, nil
}

func (uc *aiUsecase) EnhanceContent(userID string, content string) (*entities.Chat, error) {
	// Validate user exists
	_, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Enhance content using AI service
	enhancedContent, tokens, err := uc.aiService.EnhanceContent(content)
	if err != nil {
		return nil, err
	}

	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Create chat entity
	chat := &entities.Chat{
		UserID:    userObjectID,
		Request:   content,
		Response:  enhancedContent,
		Tokens:    tokens,
		CreatedAt: time.Now(),
	}

	// Save to database
	savedChat, err := uc.chatRepository.CreateChat(chat)
	if err != nil {
		return nil, err
	}

	return savedChat, nil
}

func (uc *aiUsecase) GetChatHistory(userID string) ([]entities.Chat, error) {
	// Validate user exists
	_, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get chat history from database
	chats, err := uc.chatRepository.GetChatsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (uc *aiUsecase) DeleteChat(userID string, chatID string) error {
	// Validate user exists
	_, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Get chat to check ownership
	chat, err := uc.chatRepository.GetChatByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	// Check if user owns this chat
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	if chat.UserID != userObjectID {
		return errors.New("you can only delete your own chats")
	}

	// Delete chat
	return uc.chatRepository.DeleteChat(chatID)
}
