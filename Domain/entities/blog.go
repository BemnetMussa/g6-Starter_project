package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog represents a blog post document in MongoDB.
type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID  primitive.ObjectID `bson:"author_id" json:"author_id"` // ref users._id
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Tags      []string           `bson:"tags" json:"tags"`
	ViewCount int                `bson:"view_count" json:"view_count"`
	Likes     int64              `bson:"likes,omitempty" json:"likes,omitempty"`
	Dislikes  int64              `bson:"dislikes,omitempty" json:"dislikes,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
