package usecases

import (
	"context"
	"errors"
	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/mongodb/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. DEFINE THE INTERFACE (The Contract)
// This lists all the functions our usecase must have.
type ICommentUsecase interface {
	CreateComment(ctx context.Context, blogIDStr string, authorID primitive.ObjectID, content string) (*entities.Comment, error)
}

// 2. DEFINE THE STRUCT (The "Employee" that does the work)
// This is the part that was missing.
// It holds the database repositories that it needs to do its job.
type CommentUsecase struct {
	commentRepo repositories.ICommentRepository
	blogRepo    repositories.IBlogRepository
}

// 3. DEFINE THE CONSTRUCTOR (The "Hiring" function)
// This is how we will create a new CommentUsecase from main.go.
// It takes the dependencies (the repositories) and returns a new usecase.
func NewCommentUsecase(commentRepo repositories.ICommentRepository, blogRepo repositories.IBlogRepository) ICommentUsecase {
	return &CommentUsecase{
		commentRepo: commentRepo,
		blogRepo:    blogRepo,
	}
}

// 4. YOUR EXISTING METHOD
// This method is attached to the CommentUsecase struct and now it will be valid.
func (uc *CommentUsecase) CreateComment(ctx context.Context, blogIDStr string, authorID primitive.ObjectID, content string) (*entities.Comment, error) {
	if content == "" {
		return nil, errors.New("comment content cannot be empty")
	}
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		return nil, errors.New("invalid blog ID format")
	}
	comment := &entities.Comment{
		BlogID:    blogID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	createdComment, err := uc.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	go uc.blogRepo.IncrementCommentCount(context.Background(), blogID)

	return createdComment, nil
}
