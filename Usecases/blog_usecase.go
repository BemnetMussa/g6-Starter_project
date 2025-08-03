package usecases

import (
	"context"
	"errors"
	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/mongodb/repositories"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IBlogUsecase defines the logic for blog posts
type IBlogUsecase interface {
	// CRUD usecases
	CreatePost(ctx context.Context, blog *entities.Blog, authorID primitive.ObjectID) (*entities.Blog, error)
	GetPostByID(ctx context.Context, postID string, requestingUserID *primitive.ObjectID) (*entities.Blog, error)
	UpdatePost(ctx context.Context, postID string, updateData *entities.Blog, requestingUserID primitive.ObjectID) (*entities.Blog, error)
	DeletePost(ctx context.Context, postID string, requestingUserID primitive.ObjectID, requestingUserRole string) error
	// filter & Search usecases
	ListPosts(ctx context.Context, tag, authorName, title, sortBy string, startDate, endDate *time.Time, page, limit int64) ([]entities.Blog, int64, error)
	// Popularity usecases
	LikePost(ctx context.Context, postID string, userID primitive.ObjectID) error
	DislikePost(ctx context.Context, postID string, userID primitive.ObjectID) error
}

type blogUsecase struct {
	blogRepo        repositories.IBlogRepository
	interactionRepo repositories.IBlogInteractionRepository
	userRepo        entities.UserRepository
}

// NewBlogUsecase creates a new blog usecase instance
func 	NewBlogUsecase(
			blogRepo repositories.IBlogRepository, 
			interactionRepo repositories.IBlogInteractionRepository, 
			userRepo entities.UserRepository) IBlogUsecase {

	return &blogUsecase{
		blogRepo:        blogRepo,
		interactionRepo: interactionRepo,
		userRepo:        userRepo,
	}
}

// CreatePost creates a new blog post with author and timestamps
func (uc *blogUsecase) CreatePost(ctx context.Context, post *entities.Blog, authorID primitive.ObjectID) (*entities.Blog, error) {
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.AuthorID = authorID
	return uc.blogRepo.Create(ctx, post)
}

// GetPostByID retrieves a blog post by ID and tracks user view
func (uc *blogUsecase) GetPostByID(ctx context.Context, postID string, userID *primitive.ObjectID) (*entities.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	post, err := uc.blogRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if userID != nil {
		interaction := &entities.BlogInteraction{
			BlogID: objectID,
			UserID: *userID,
			Viewed: true,
		}
		_ = uc.interactionRepo.Upsert(ctx, interaction)
	}

	likes, dislikes, views, err := uc.interactionRepo.GetPopularityCounts(ctx, objectID)
	if err != nil {
		return post, nil
	}

	post.Likes = int(likes)
	post.Dislikes = int(dislikes)
	post.ViewCount = int(views)

	return post, nil
}

// UpdatePost updates blog content, tags, and timestamps if user is the author
func (uc *blogUsecase) UpdatePost(ctx context.Context, postID string, updateData *entities.Blog, requestingUserID primitive.ObjectID) (*entities.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	originalPost, err := uc.blogRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if originalPost.AuthorID != requestingUserID {
		return nil, errors.New("forbidden: you are not the author of this post")
	}

	originalPost.Title = updateData.Title
	originalPost.Content = updateData.Content
	originalPost.Tags = updateData.Tags
	originalPost.UpdatedAt = time.Now()

	err = uc.blogRepo.Update(ctx, originalPost)
	if err != nil {
		return nil, err
	}
	return originalPost, nil
}

// DeletePost deletes a blog post if the requester is the author or an admin
func (uc *blogUsecase) DeletePost(ctx context.Context, postID string, requestingUserID primitive.ObjectID, requestingUserRole string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New("invalid post ID format")
	}

	postToDelete, err := uc.blogRepo.FindByID(ctx, objectID)
	if err != nil {
		return errors.New("post not found")
	}

	isOwner := postToDelete.AuthorID == requestingUserID
	isAdmin := requestingUserRole == "admin"

	if !isOwner && !isAdmin {
		return errors.New("forbidden: you are not authorized to delete this post")
	}

	return uc.blogRepo.Delete(ctx, objectID)
}

// ListPosts returns filtered and paginated blog posts with search and sort options
func (uc *blogUsecase) ListPosts(
    ctx context.Context,
    tag string,
    authorIDStr string,
    title string,
    sortBy string,
    startDate, endDate *time.Time,
    page, limit int64,
) ([]entities.Blog, int64, error) {
    var authorID *primitive.ObjectID
    var tags []string

    // Convert authorID string to ObjectID pointer 
    if authorIDStr != "" {
        id, err := primitive.ObjectIDFromHex(authorIDStr)
        if err != nil {
            return nil, 0, errors.New("invalid author ID")
        }
        authorID = &id
    }

    if tag != "" {
        tags = strings.Split(tag, ",")
    }

    options := repositories.SearchFilterOptions{
        AuthorID:  authorID,
        Tags:      tags,
        Title:     title,
        Page:      page,
        Limit:     limit,
        StartDate: startDate,
        EndDate:   endDate,
        SortBy:    sortBy,
    }

    return uc.blogRepo.Find(ctx, options)
}

// LikePost records a like interaction for a blog post by a user
func (uc *blogUsecase) LikePost(ctx context.Context, postID string, userID primitive.ObjectID) error {
	blogObjectID, err := primitive.ObjectIDFromHex(postID)
	reactionLike := "like"
	if err != nil {
		return errors.New("invalid post ID format")
	}
	interaction := &entities.BlogInteraction{
		BlogID:   blogObjectID,
		UserID:   userID,
		Reaction: &reactionLike,
	}
	return uc.interactionRepo.Upsert(ctx, interaction)
}

// DislikePost records a dislike interaction for a blog post by a user
func (uc *blogUsecase) DislikePost(ctx context.Context, postID string, userID primitive.ObjectID) error {
	blogObjectID, err := primitive.ObjectIDFromHex(postID)
	reactionDislike := "dislike"
	if err != nil {
		return errors.New("invalid post ID format")
	}
	interaction := &entities.BlogInteraction{
		BlogID:   blogObjectID,
		UserID:   userID,
		Reaction: &reactionDislike,
	}
	return uc.interactionRepo.Upsert(ctx, interaction)
}
