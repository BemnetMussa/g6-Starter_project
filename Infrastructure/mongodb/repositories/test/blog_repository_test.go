package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/mongodb/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogTestSuite struct {
	client           *mongo.Client
	database         *mongo.Database
	blogCollection   *mongo.Collection
	interactionCollection *mongo.Collection
	commentCollection *mongo.Collection
	blogRepo         repositories.IBlogRepository
	interactionRepo  repositories.IBlogInteractionRepository
	commentRepo      repositories.ICommentRepository
	config           *TestConfig
}

func setupBlogTestSuite(t *testing.T) *BlogTestSuite {
	config := GetTestConfig()
	// Use a unique database name for blog tests to avoid conflicts
	config.DatabaseName = fmt.Sprintf("test_blog_api_%d", time.Now().Unix())
	client, database, _ := SetupTestDatabase(t, config)
	
	// Create collections for blog testing - use the same collection names as the repository
	blogCollection := database.Collection("blogs")
	interactionCollection := database.Collection("blog_interactions")
	commentCollection := database.Collection("comments")

	// Clear collections before each test
	_, err := blogCollection.DeleteMany(context.TODO(), bson.M{})
	require.NoError(t, err)
	_, err = interactionCollection.DeleteMany(context.TODO(), bson.M{})
	require.NoError(t, err)
	_, err = commentCollection.DeleteMany(context.TODO(), bson.M{})
	require.NoError(t, err)

	// Create repositories
	blogRepo := repositories.NewBlogRepository(database)
	interactionRepo := repositories.NewBlogInteractionRepository(database)
	commentRepo := repositories.NewCommentRepository(database)

	return &BlogTestSuite{
		client:             client,
		database:           database,
		blogCollection:     blogCollection,
		interactionCollection: interactionCollection,
		commentCollection:  commentCollection,
		blogRepo:           blogRepo,
		interactionRepo:    interactionRepo,
		commentRepo:        commentRepo,
		config:             config,
	}
}

func (ts *BlogTestSuite) teardown(t *testing.T) {
	CleanupTestDatabase(t, ts.client, ts.database)
}

func createTestBlog(authorID primitive.ObjectID) *entities.Blog {
	now := time.Now()
	return &entities.Blog{
		AuthorID:     authorID,
		Title:        "Test Blog Post",
		Content:      "This is a test blog post content.",
		Tags:         []string{"test", "blog", "go"},
		ViewCount:    0,
		Likes:        0,
		Dislikes:     0,
		CommentCount: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func createTestBlogWithCustomFields(authorID primitive.ObjectID, title, content string, tags []string) *entities.Blog {
	blog := createTestBlog(authorID)
	blog.Title = title
	blog.Content = content
	blog.Tags = tags
	return blog
}

func createTestBlogInteraction(blogID, userID primitive.ObjectID) *entities.BlogInteraction {
	now := time.Now()
	return &entities.BlogInteraction{
		BlogID:       blogID,
		UserID:       userID,
		Reaction:     nil,
		Viewed:       true,
		InteractedAt: now,
	}
}

func createTestComment(blogID, authorID primitive.ObjectID) *entities.Comment {
	now := time.Now()
	return &entities.Comment{
		BlogID:    blogID,
		AuthorID:  authorID,
		Content:   "This is a test comment.",
		CreatedAt: now,
	}
}

func TestBlogRepository_Create(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should create blog successfully", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		blog := createTestBlog(authorID)

		createdBlog, err := ts.blogRepo.Create(context.TODO(), blog)

		assert.NoError(t, err)
		assert.NotNil(t, createdBlog)
		assert.NotEmpty(t, createdBlog.ID)
		assert.Equal(t, authorID, createdBlog.AuthorID)
		assert.Equal(t, blog.Title, createdBlog.Title)
		assert.Equal(t, blog.Content, createdBlog.Content)
		assert.Equal(t, blog.Tags, createdBlog.Tags)
		assert.Equal(t, 0, createdBlog.ViewCount)
		assert.Equal(t, 0, createdBlog.Likes)
		assert.Equal(t, 0, createdBlog.Dislikes)
		assert.Equal(t, 0, createdBlog.CommentCount)
	})

	t.Run("should create blog with custom fields", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		title := "Custom Blog Title"
		content := "Custom blog content with more details."
		tags := []string{"custom", "test", "example"}
		
		blog := createTestBlogWithCustomFields(authorID, title, content, tags)

		createdBlog, err := ts.blogRepo.Create(context.TODO(), blog)

		assert.NoError(t, err)
		assert.NotNil(t, createdBlog)
		assert.Equal(t, title, createdBlog.Title)
		assert.Equal(t, content, createdBlog.Content)
		assert.Equal(t, tags, createdBlog.Tags)
	})
}

func TestBlogRepository_FindByID(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should find blog by ID", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		blog := createTestBlog(authorID)
		
		createdBlog, err := ts.blogRepo.Create(context.TODO(), blog)
		require.NoError(t, err)

		foundBlog, err := ts.blogRepo.FindByID(context.TODO(), createdBlog.ID)

		assert.NoError(t, err)
		assert.NotNil(t, foundBlog)
		assert.Equal(t, createdBlog.ID, foundBlog.ID)
		assert.Equal(t, createdBlog.Title, foundBlog.Title)
		assert.Equal(t, createdBlog.Content, foundBlog.Content)
	})

	t.Run("should return error for non-existent blog", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		
		_, err := ts.blogRepo.FindByID(context.TODO(), nonExistentID)

		assert.Error(t, err)
	})
}

func TestBlogRepository_Update(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should update blog successfully", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		blog := createTestBlog(authorID)
		
		createdBlog, err := ts.blogRepo.Create(context.TODO(), blog)
		require.NoError(t, err)

		// Update blog
		createdBlog.Title = "Updated Blog Title"
		createdBlog.Content = "Updated blog content."
		createdBlog.Tags = []string{"updated", "test"}

		err = ts.blogRepo.Update(context.TODO(), createdBlog)

		assert.NoError(t, err)

		// Verify update
		updatedBlog, err := ts.blogRepo.FindByID(context.TODO(), createdBlog.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Blog Title", updatedBlog.Title)
		assert.Equal(t, "Updated blog content.", updatedBlog.Content)
		assert.Equal(t, []string{"updated", "test"}, updatedBlog.Tags)
	})

	t.Run("should return error for non-existent blog", func(t *testing.T) {
		blog := createTestBlog(primitive.NewObjectID())
		blog.ID = primitive.NewObjectID()

		err := ts.blogRepo.Update(context.TODO(), blog)

		assert.Error(t, err)
	})
}

func TestBlogRepository_Delete(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should delete blog successfully", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		blog := createTestBlog(authorID)
		
		createdBlog, err := ts.blogRepo.Create(context.TODO(), blog)
		require.NoError(t, err)

		err = ts.blogRepo.Delete(context.TODO(), createdBlog.ID)

		assert.NoError(t, err)

		// Verify blog is deleted
		_, err = ts.blogRepo.FindByID(context.TODO(), createdBlog.ID)
		assert.Error(t, err)
	})
}

func TestBlogRepository_Find(t *testing.T) {
	t.Run("should find blogs with pagination", func(t *testing.T) {
		ts := setupBlogTestSuite(t)
		defer ts.teardown(t)

		authorID := primitive.NewObjectID()
		
		// Create multiple blogs
		for i := 1; i <= 5; i++ {
			blog := createTestBlogWithCustomFields(authorID, fmt.Sprintf("Blog %d", i), fmt.Sprintf("Content %d", i), []string{"test"})
			_, err := ts.blogRepo.Create(context.TODO(), blog)
			require.NoError(t, err)
		}

		options := repositories.SearchFilterOptions{
			Page:  1,
			Limit: 3,
		}

		blogs, count, err := ts.blogRepo.Find(context.TODO(), options)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.Len(t, blogs, 3)
	})

	t.Run("should find blogs by author", func(t *testing.T) {
		ts := setupBlogTestSuite(t)
		defer ts.teardown(t)

		authorID := primitive.NewObjectID()
		otherAuthorID := primitive.NewObjectID()
		
		blog1 := createTestBlogWithCustomFields(authorID, "Blog 1", "Content 1", []string{"test"})
		blog2 := createTestBlogWithCustomFields(authorID, "Blog 2", "Content 2", []string{"test"})
		blog3 := createTestBlogWithCustomFields(otherAuthorID, "Blog 3", "Content 3", []string{"test"})

		_, err := ts.blogRepo.Create(context.TODO(), blog1)
		require.NoError(t, err)
		_, err = ts.blogRepo.Create(context.TODO(), blog2)
		require.NoError(t, err)
		_, err = ts.blogRepo.Create(context.TODO(), blog3)
		require.NoError(t, err)

		options := repositories.SearchFilterOptions{
			AuthorID: &authorID,
		}

		blogs, count, err := ts.blogRepo.Find(context.TODO(), options)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
		assert.Len(t, blogs, 2)
	})

	t.Run("should find blogs by tags", func(t *testing.T) {
		ts := setupBlogTestSuite(t)
		defer ts.teardown(t)

		authorID := primitive.NewObjectID()
		
		blog1 := createTestBlogWithCustomFields(authorID, "Blog 1", "Content 1", []string{"go", "test"})
		blog2 := createTestBlogWithCustomFields(authorID, "Blog 2", "Content 2", []string{"python", "test"})
		blog3 := createTestBlogWithCustomFields(authorID, "Blog 3", "Content 3", []string{"javascript"})

		_, err := ts.blogRepo.Create(context.TODO(), blog1)
		require.NoError(t, err)
		_, err = ts.blogRepo.Create(context.TODO(), blog2)
		require.NoError(t, err)
		_, err = ts.blogRepo.Create(context.TODO(), blog3)
		require.NoError(t, err)

		options := repositories.SearchFilterOptions{
			Tags: []string{"test"},
		}

		blogs, count, err := ts.blogRepo.Find(context.TODO(), options)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
		assert.Len(t, blogs, 2)
	})

	t.Run("should find blogs by title", func(t *testing.T) {
		ts := setupBlogTestSuite(t)
		defer ts.teardown(t)

		authorID := primitive.NewObjectID()
		
		blog1 := createTestBlogWithCustomFields(authorID, "Go Programming", "Content 1", []string{"go"})
		blog2 := createTestBlogWithCustomFields(authorID, "Python Programming", "Content 2", []string{"python"})

		_, err := ts.blogRepo.Create(context.TODO(), blog1)
		require.NoError(t, err)
		_, err = ts.blogRepo.Create(context.TODO(), blog2)
		require.NoError(t, err)

		options := repositories.SearchFilterOptions{
			Title: "Go",
		}

		blogs, count, err := ts.blogRepo.Find(context.TODO(), options)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
		assert.Len(t, blogs, 1)
		assert.Equal(t, "Go Programming", blogs[0].Title)
	})
}

func TestBlogRepository_UpdateCounts(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should update blog counts", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		blog := createTestBlog(authorID)
		
		createdBlog, err := ts.blogRepo.Create(context.TODO(), blog)
		require.NoError(t, err)

		err = ts.blogRepo.UpdateCounts(context.TODO(), createdBlog.ID, 10, 2)

		assert.NoError(t, err)

		// Verify counts were updated
		updatedBlog, err := ts.blogRepo.FindByID(context.TODO(), createdBlog.ID)
		assert.NoError(t, err)
		assert.Equal(t, 10, updatedBlog.Likes)
		assert.Equal(t, 2, updatedBlog.Dislikes)
	})
}

func TestBlogRepository_IncrementCommentCount(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should increment comment count", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		blog := createTestBlog(authorID)
		
		createdBlog, err := ts.blogRepo.Create(context.TODO(), blog)
		require.NoError(t, err)

		err = ts.blogRepo.IncrementCommentCount(context.TODO(), createdBlog.ID)

		assert.NoError(t, err)

		// Verify comment count was incremented
		updatedBlog, err := ts.blogRepo.FindByID(context.TODO(), createdBlog.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, updatedBlog.CommentCount)
	})
}

func TestBlogInteractionRepository_Upsert(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should create new interaction", func(t *testing.T) {
		blogID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		interaction := createTestBlogInteraction(blogID, userID)

		err := ts.interactionRepo.Upsert(context.TODO(), interaction)

		assert.NoError(t, err)
	})

	t.Run("should update existing interaction", func(t *testing.T) {
		blogID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		interaction := createTestBlogInteraction(blogID, userID)

		// Create interaction
		err := ts.interactionRepo.Upsert(context.TODO(), interaction)
		require.NoError(t, err)

		// Update interaction
		likeReaction := "like"
		interaction.Reaction = &likeReaction
		err = ts.interactionRepo.Upsert(context.TODO(), interaction)

		assert.NoError(t, err)

		// Verify interaction was updated
		foundInteraction, err := ts.interactionRepo.FindByBlogAndUser(context.TODO(), blogID, userID)
		assert.NoError(t, err)
		assert.Equal(t, "like", *foundInteraction.Reaction)
	})
}

func TestBlogInteractionRepository_GetPopularityCounts(t *testing.T) {
	t.Run("should get popularity counts", func(t *testing.T) {
		ts := setupBlogTestSuite(t)
		defer ts.teardown(t)

		blogID := primitive.NewObjectID()
		user1ID := primitive.NewObjectID()
		user2ID := primitive.NewObjectID()
		user3ID := primitive.NewObjectID()

		// Create interactions
		interaction1 := createTestBlogInteraction(blogID, user1ID)
		likeReaction := "like"
		interaction1.Reaction = &likeReaction
		err := ts.interactionRepo.Upsert(context.TODO(), interaction1)
		require.NoError(t, err)

		interaction2 := createTestBlogInteraction(blogID, user2ID)
		dislikeReaction := "dislike"
		interaction2.Reaction = &dislikeReaction
		err = ts.interactionRepo.Upsert(context.TODO(), interaction2)
		require.NoError(t, err)

		interaction3 := createTestBlogInteraction(blogID, user3ID)
		interaction3.Viewed = true
		err = ts.interactionRepo.Upsert(context.TODO(), interaction3)
		require.NoError(t, err)

		likes, dislikes, views, err := ts.interactionRepo.GetPopularityCounts(context.TODO(), blogID)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), likes)
		assert.Equal(t, int64(1), dislikes)
		assert.Equal(t, int64(1), views)
	})
}

func TestBlogInteractionRepository_FindByBlogAndUser(t *testing.T) {
	t.Run("should find interaction by blog and user", func(t *testing.T) {
		ts := setupBlogTestSuite(t)
		defer ts.teardown(t)

		blogID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		interaction := createTestBlogInteraction(blogID, userID)

		err := ts.interactionRepo.Upsert(context.TODO(), interaction)
		require.NoError(t, err)

		foundInteraction, err := ts.interactionRepo.FindByBlogAndUser(context.TODO(), blogID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, foundInteraction)
		assert.Equal(t, blogID, foundInteraction.BlogID)
		assert.Equal(t, userID, foundInteraction.UserID)
	})

	t.Run("should return error for non-existent interaction", func(t *testing.T) {
		ts := setupBlogTestSuite(t)
		defer ts.teardown(t)

		blogID := primitive.NewObjectID()
		userID := primitive.NewObjectID()

		_, err := ts.interactionRepo.FindByBlogAndUser(context.TODO(), blogID, userID)

		assert.Error(t, err)
	})
}

func TestCommentRepository_Create(t *testing.T) {
	ts := setupBlogTestSuite(t)
	defer ts.teardown(t)

	t.Run("should create comment successfully", func(t *testing.T) {
		blogID := primitive.NewObjectID()
		authorID := primitive.NewObjectID()
		comment := createTestComment(blogID, authorID)

		createdComment, err := ts.commentRepo.Create(context.TODO(), comment)

		assert.NoError(t, err)
		assert.NotNil(t, createdComment)
		assert.NotEmpty(t, createdComment.ID)
		assert.Equal(t, blogID, createdComment.BlogID)
		assert.Equal(t, authorID, createdComment.AuthorID)
		assert.Equal(t, comment.Content, createdComment.Content)
	})
} 