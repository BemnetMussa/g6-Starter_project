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
func NewBlogUsecase(
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

	post.Likes = int64(likes)
	post.Dislikes = int64(dislikes)
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

// new
// --- Core Like/Dislike Logic ---

const (
	reactionLike    = "like"
	reactionDislike = "dislike"
)

// LikePost handles the logic for a user liking a post.
func (uc *blogUsecase) LikePost(ctx context.Context, blogIDStr string, userID primitive.ObjectID) error {
	return uc.reactToPost(ctx, blogIDStr, userID, reactionLike)
}

// DislikePost handles the logic for a user disliking a post.
func (uc *blogUsecase) DislikePost(ctx context.Context, blogIDStr string, userID primitive.ObjectID) error {
	return uc.reactToPost(ctx, blogIDStr, userID, reactionDislike)
}

// reactToPost is the shared private method that contains the core logic.
func (uc *blogUsecase) reactToPost(ctx context.Context, blogIDStr string, userID primitive.ObjectID, newReactionType string) error {
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		return errors.New("invalid blog ID format")
	}

	// 1. Find if a previous interaction exists from this user for this blog.
	existingInteraction, err := uc.interactionRepo.FindByBlogAndUser(ctx, blogID, userID)
	if err != nil {
		return err
	}

	var finalReaction *string

	// 2. Determine the new state of the reaction based on the rules.
	if existingInteraction != nil && existingInteraction.Reaction != nil && *existingInteraction.Reaction == newReactionType {
		// User is clicking the same button again (e.g., clicking "like" when it's already liked).
		// This means we toggle it OFF.
		finalReaction = nil // Setting reaction to nil means "no reaction"
	} else {
		// User is either reacting for the first time, or changing their reaction (e.g., from dislike to like).
		finalReaction = &newReactionType
	}

	// 3. Create the interaction object and upsert it into the database.
	interaction := &entities.BlogInteraction{
		BlogID:   blogID,
		UserID:   userID,
		Reaction: finalReaction,
		// If you track views, you'd handle the logic for setting `Viewed` here as well.
		// For a like/dislike action, we can assume it's also a view.
		Viewed: true,
	}
	if err := uc.interactionRepo.Upsert(ctx, interaction); err != nil {
		return err
	}

	// 4. Get the new total counts of likes and dislikes for the post.
	likes, dislikes, _, err := uc.interactionRepo.GetPopularityCounts(ctx, blogID)
	if err != nil {
		return err
	}

	// 5. Update the denormalized counts on the main Blog document.
	return uc.blogRepo.UpdateCounts(ctx, blogID, likes, dislikes)
}
