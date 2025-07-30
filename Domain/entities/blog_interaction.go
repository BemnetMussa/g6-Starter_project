package entities

import (
	"time"
)

// BlogInteraction tracks user interactions with blogs.
type BlogInteraction struct {
    ID           string    `bson:"_id,omitempty" json:"id"`
    BlogID       string    `bson:"blog_id" json:"blog_id"`   // ref blogs._id
    UserID       string    `bson:"user_id" json:"user_id"`   // ref users._id
    Reaction     *string   `bson:"reaction,omitempty" json:"reaction,omitempty"` // "like", "dislike", or nil
    Viewed       bool      `bson:"viewed" json:"viewed"`
    InteractedAt time.Time `bson:"interacted_at" json:"interacted_at"`
}
