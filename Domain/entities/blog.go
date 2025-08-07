package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog represents a blog post document in MongoDB.
type Blog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID     primitive.ObjectID `bson:"author_id" json:"author_id"` // ref users._id
	Title        string             `bson:"title" json:"title" binding:"required"`
	Content      string             `bson:"content" json:"content" binding:"required"`
	Tags         []string           `bson:"tags" json:"tags"`
	ViewCount    int                `bson:"view_count" json:"view_count"`
	Likes        int                `bson:"likes" json:"likes"`
	Dislikes     int                `bson:"dislikes" json:"dislikes"`
	CommentCount int                `bson:"comment_count" json:"comment_count"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
