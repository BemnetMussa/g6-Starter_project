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
	userRepo        repositories.IUserRepository // IUserRepository is the userrepository interface
}

// NewBlogUsecase is the constructor.
func NewBlogUsecase(br repositories.IBlogRepository, ir repositories.IBlogInteractionRepository, userRepo repositories.IUserRepository) IBlogUsecase {
	return &blogUsecase{
		blogRepo:        br,
		interactionRepo: ir,
		userRepo:        ur,
	}
}

//  Method Implementations

func (uc *blogUsecase) CreatePost(ctx context.Context, post *entities.Blog, authorID primitive.ObjectID) (*entities.Blog, error) {
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.AuthorID = authorID
	return uc.blogRepo.Create(ctx, post)
}

func (uc *blogUsecase) GetPostByID(ctx context.Context, postID string, userID *primitive.ObjectID) (*entities.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	// fetch the blog post
	post, err := uc.blogRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// track the view interaction.
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

func (uc *blogUsecase) UpdatePost(ctx context.Context, postID string, updateData *entities.Blog, requestingUserID primitive.ObjectID) (*entities.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	// Get the original post
	originalPost, err := uc.blogRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Check if the requester is the author
	if originalPost.AuthorID != requestingUserID {
		return nil, errors.New("forbidden: you are not the author of this post")
	}

	// Apply the updates
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

func (uc *blogUsecase) DeletePost(ctx context.Context, postID string, requestingUserID primitive.ObjectID, requestingUserRole string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New("invalid post ID format")
	}

	// Fetch the post
	postToDelete, err := uc.blogRepo.FindByID(ctx, objectID)
	if err != nil {
		return errors.New("post not found")
	}

	// Check if the requester is the author OR an admin.
	isOwner := postToDelete.AuthorID == requestingUserID
	isAdmin := requestingUserRole == "admin"

	if !isOwner && !isAdmin {
		return errors.New("forbidden: you are not authorized to delete this post")
	}

	return uc.blogRepo.Delete(ctx, objectID)
}

func (uc *blogUsecase) ListPosts(ctx context.Context, tag, authorName, title, sortBy string, startDate, endDate *time.Time, page, limit int64) ([]entities.Blog, int64, error) {
	var authorID *primitive.ObjectID
	var tags []string

	if authorName != "" {
		author, err := uc.userRepo.FindByUsername(ctx, authorName)
		if err == nil {
			authorID = &author.ID
		} else {
			return []entities.Blog{}, 0, nil
		}
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
