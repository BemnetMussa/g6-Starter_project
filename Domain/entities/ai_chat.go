package entities

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Chat represents an AI chat interaction document in MongoDB.
type Chat struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
    Request   string             `bson:"request" json:"request"`
    Response  string             `bson:"response" json:"response"`
    Tokens    int                `bson:"tokens" json:"tokens"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// interface for repository to use
type ChatRepository interface {
	CreateChat(chat *Chat) (*Chat, error)
	GetChatByID(id string) (*Chat, error)
	GetChatsByUserID(userID string) ([]Chat, error)
	DeleteChat(id string) error
}