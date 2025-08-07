package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BlogID    primitive.ObjectID `bson:"blog_id" json:"blog_id"`     // Which blog it's for
	AuthorID  primitive.ObjectID `bson:"author_id" json:"author_id"` // Who wrote it
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
