package repositories

import (
	"context"
	"errors"
	// "time"

	"g6_starter_project/Domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepositoryImpl struct {
	db *mongo.Collection
}

func NewChatRepository(db *mongo.Collection) entities.ChatRepository {
	return &ChatRepositoryImpl{db: db}
}

func (r *ChatRepositoryImpl) CreateChat(chat *entities.Chat) (*entities.Chat, error) {
	result, err := r.db.InsertOne(context.TODO(), chat)
	if err != nil {
		return nil, err
	}
	
	// Set the generated ID back to the chat object
	chat.ID = result.InsertedID.(primitive.ObjectID)
	return chat, nil
}

func (r *ChatRepositoryImpl) GetChatByID(id string) (*entities.Chat, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid chat ID")
	}

	filter := bson.M{"_id": objectID}
	var chat entities.Chat
	err = r.db.FindOne(context.TODO(), filter).Decode(&chat)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("chat not found")
		}
		return nil, err
	}
	return &chat, nil
}

func (r *ChatRepositoryImpl) GetChatsByUserID(userID string) ([]entities.Chat, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	filter := bson.M{"user_id": objectID}
	cursor, err := r.db.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var chats []entities.Chat
	if err = cursor.All(context.TODO(), &chats); err != nil {
		return nil, err
	}
	
	// Return empty slice instead of nil if no chats found
	if chats == nil {
		chats = []entities.Chat{}
	}
	
	return chats, nil
}

func (r *ChatRepositoryImpl) DeleteChat(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid chat ID")
	}

	filter := bson.M{"_id": objectID}
	_, err = r.db.DeleteOne(context.TODO(), filter)
	return err
}	
