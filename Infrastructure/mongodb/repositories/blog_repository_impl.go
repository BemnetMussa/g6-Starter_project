package repositories

import (
	"context"
	"starter_project/Domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchFilterOptions is a struct to hold all possible criteria for searching/filtering.
type SearchFilterOptions struct {
	AuthorID *primitive.ObjectID
	Tags     []string
	Title    string
	Page     int64
	Limit    int64
}

// IBlogRepository defines the contract for all blog data operations.
type IBlogRepository interface {
	// CRUD
	Create(ctx context.Context, post *entities.Blog) (*entities.Blog, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Blog, error)
	Update(ctx context.Context, post *entities.Blog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	// Advanced queries
	Find(ctx context.Context, options SearchFilterOptions) ([]entities.Blog, int64, error)
}
