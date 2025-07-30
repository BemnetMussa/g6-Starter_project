package entities

import "time"

// Blog represents a blog post document in MongoDB.
type Blog struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    AuthorID  string    `bson:"author_id" json:"author_id"` // ref users._id
    Title     string    `bson:"title" json:"title"`
    Content   string    `bson:"content" json:"content"`
    Tags      []string  `bson:"tags" json:"tags"`
    ViewCount int       `bson:"view_count" json:"view_count"`
    Likes     *int      `bson:"likes,omitempty" json:"likes,omitempty"`
    Dislikes  *int      `bson:"dislikes,omitempty" json:"dislikes,omitempty"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
