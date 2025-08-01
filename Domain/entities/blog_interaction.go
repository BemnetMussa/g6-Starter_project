package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogInteraction tracks user interactions with blogs.
type BlogInteraction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BlogID       primitive.ObjectID `bson:"blog_id" json:"blog_id"`                       // ref blogs._id
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`                       // ref users._id
	Reaction     *string            `bson:"reaction,omitempty" json:"reaction,omitempty"` // "like", "dislike", or nil
	Viewed       bool               `bson:"viewed" json:"viewed"`
	InteractedAt time.Time          `bson:"interacted_at" json:"interacted_at"`
}
