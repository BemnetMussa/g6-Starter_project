package repositories

import (
	"context"
	"errors"
	"starter_project/Domain/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SearchFilterOptions is a struct to hold all possible criteria for searching/filtering.
type SearchFilterOptions struct {
	AuthorID  *primitive.ObjectID
	Tags      []string
	Title     string
	Page      int64
	Limit     int64
	StartDate *time.Time
	EndDate   *time.Time
	SortBy    string
}

// IBlogRepository defines the contract for all blog data operations.
type IBlogRepository interface {
	// CRUD
	Create(ctx context.Context, blog *entities.Blog) (*entities.Blog, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Blog, error)
	Update(ctx context.Context, blog *entities.Blog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	// Advanced queries: seraching, filtering
	Find(ctx context.Context, options SearchFilterOptions) ([]entities.Blog, int64, error)
}

// IBlogInteractionRepository defines the contract for interaction data.
type IBlogInteractionRepository interface {
	// Upsert means "update if exists, insert if not".
	Upsert(ctx context.Context, interaction *entities.BlogInteraction) error
	GetPopularityCounts(ctx context.Context, blogID primitive.ObjectID) (likes int64, dislikes int64, views int64, err error)
}

type mongoBlogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) IBlogRepository {
	return &mongoBlogRepository{collection: db.Collection("blogs")}
}

type mongoBlogInteractionRepository struct {
	collection *mongo.Collection
}

func NewBlogInteractionRepository(db *mongo.Database) IBlogInteractionRepository {
	return &mongoBlogInteractionRepository{collection: db.Collection("blog_interactions")}
}

// // --- Method Implementations ---

func (r *mongoBlogRepository) Create(ctx context.Context, blog *entities.Blog) (*entities.Blog, error) {
	result, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, err
	}
	blog.ID = result.InsertedID.(primitive.ObjectID)
	return blog, nil
}

func (r *mongoBlogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Blog, error) {
	var blog entities.Blog
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *mongoBlogRepository) Update(ctx context.Context, blog *entities.Blog) error {
	filter := bson.M{"_id": blog.ID}
	update := bson.M{"$set": blog}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err == nil && result.MatchedCount == 0 {
		return errors.New("post not found to update")
	}
	return err
}

func (r *mongoBlogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err == nil && result.DeletedCount == 0 {
		return errors.New("post not found to delete")
	}
	return err
}

func (r *mongoBlogRepository) Find(ctx context.Context, filterOptions SearchFilterOptions) ([]entities.Blog, int64, error) {
	filter := bson.M{}
	if filterOptions.AuthorID != nil {
		filter["author_id"] = filterOptions.AuthorID
	}
	if filterOptions.Title != "" {
		// regex for case-insensitive "contains" search
		filter["title"] = bson.M{"$regex": filterOptions.Title, "$options": "i"}
	}
	if len(filterOptions.Tags) > 0 {
		// $in to match any of the tags in the array
		filter["tags"] = bson.M{"$in": filterOptions.Tags}
	}

	if filterOptions.StartDate != nil && filterOptions.EndDate != nil {
		filter["created_at"] = bson.M{
			"$gte": filterOptions.StartDate, // gte = Greter Than or Equal to
			"$lte": filterOptions.EndDate,   // lte = Less Than or Equal to
		}
	} else if filterOptions.StartDate != nil {
		filter["created_at"] = bson.M{"$gte": filterOptions.StartDate}
	} else if filterOptions.EndDate != nil {
		filter["created_at"] = bson.M{"$lte": filterOptions.EndDate}
	}
	// pagination options
	findOptions := options.Find()
	findOptions.SetLimit(filterOptions.Limit)
	findOptions.SetSkip((filterOptions.Page - 1) * filterOptions.Limit)
	sort := bson.M{"created_at": -1}

	if filterOptions.SortBy == "popularity" {
		sort = bson.M{"likes": -1}
	}

	findOptions.SetSort(sort)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []entities.Blog
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return posts, totalCount, nil
}

// Bloginteraction - Method Implementations -

func (r *mongoBlogInteractionRepository) Upsert(ctx context.Context, interaction *entities.BlogInteraction) error {
	filter := bson.M{
		"blog_id": interaction.BlogID,
		"user_id": interaction.UserID,
	}

	update := bson.M{
		"$set": bson.M{
			"reaction":     interaction.Reaction,
			"viewed":       interaction.Viewed,
			"interactedat": time.Now(),
		},
	}

	// options.Update().SetUpsert(true) tells MongoDB to perform an upsert.
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetPopularityCounts calculates the totals for likes, dislikes, and views.
func (r *mongoBlogInteractionRepository) GetPopularityCounts(ctx context.Context, blogID primitive.ObjectID) (likes int64, dislikes int64, views int64, err error) {

	likes, err = r.collection.CountDocuments(ctx, bson.M{"blog_id": blogID, "reaction": "like"})
	if err != nil {
		return 0, 0, 0, err
	}

	dislikes, err = r.collection.CountDocuments(ctx, bson.M{"blog_id": blogID, "reaction": "dislike"})
	if err != nil {
		return 0, 0, 0, err
	}

	views, err = r.collection.CountDocuments(ctx, bson.M{"blog_id": blogID, "viewed": true})
	if err != nil {
		return 0, 0, 0, err
	}

	return likes, dislikes, views, nil
}
