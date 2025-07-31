package usecases

import (
	"context"
	"g6_starter_project/Domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IBlogUsecase defines the logic for blog posts
type IBlogUsecase interface {
	// CRUD usecases
	CreatePost(ctx context.Context, post *entities.Blog) (*entities.Blog, error)
	GetPostByID(ctx context.Context, postID string, requestingUserID primitive.ObjectID) (*entities.Blog, error)
	UpdatePost(ctx context.Context, postID string, updateData *entities.Blog, requestingUserID primitive.ObjectID) (*entities.Blog, error)
	DeletePost(ctx context.Context, postID string, requestingUserID primitive.ObjectID) error
	// Retrieval & Search usecases

	ListPosts(ctx context.Context, tag, authorName, title string, page, limit int64) ([]entities.Blog, int64, error)
	// Popularity usecases
	LikePost(ctx context.Context, postID string, userID primitive.ObjectID) error
	DislikePost(ctx context.Context, postID string, userID primitive.ObjectID) error
	Addview(ctx context.Context, postID string) error
}
