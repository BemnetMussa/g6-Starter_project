package repositories

import (
	"context"
	"errors"
	"g6_starter_project/Domain/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SearchFilterOptions is a struct to hold all possible criteria for searching/filtering.
type SearchFilterOptions struct {
	AuthorID      *primitive.ObjectID
	Tags          []string
	Title         string
	Page          int64
	Limit         int64
	StartDate     *time.Time
	EndDate       *time.Time
	SortBy        string
	MinPopularity *int64
	MaxPopularity *int64
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
	UpdateCounts(ctx context.Context, blogID primitive.ObjectID, likes, dislikes int64) error //new
	IncrementCommentCount(ctx context.Context, blogID primitive.ObjectID) error
}

// IBlogInteractionRepository defines the contract for interaction data.
type IBlogInteractionRepository interface {
	// Upsert means "update if exists, insert if not".
	Upsert(ctx context.Context, interaction *entities.BlogInteraction) error
	GetPopularityCounts(ctx context.Context, blogID primitive.ObjectID) (likes int64, dislikes int64, views int64, err error)
	FindByBlogAndUser(ctx context.Context, blogID, userID primitive.ObjectID) (*entities.BlogInteraction, error) //new
}

type ICommentRepository interface {
	Create(ctx context.Context, comment *entities.Comment) (*entities.Comment, error)
}

type mongoBlogRepository struct {
	collection *mongo.Collection
}

type mongoCommentRepository struct {
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
func NewCommentRepository(db *mongo.Database) ICommentRepository {
	return &mongoCommentRepository{collection: db.Collection("comments")}
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

// In Infrastructure/mongodb/repositories/blog_repository_impl.go

func (r *mongoBlogRepository) Find(ctx context.Context, filterOptions SearchFilterOptions) ([]entities.Blog, int64, error) {
	// --- STAGE 1: Build the Complete Filter ---
	filter := bson.M{}

	if filterOptions.AuthorID != nil {
		filter["author_id"] = filterOptions.AuthorID
	}
	if filterOptions.Title != "" {
		filter["title"] = bson.M{"$regex": filterOptions.Title, "$options": "i"}
	}
	if len(filterOptions.Tags) > 0 && filterOptions.Tags[0] != "" {
		filter["tags"] = bson.M{"$in": filterOptions.Tags}
	}
	if filterOptions.StartDate != nil && filterOptions.EndDate != nil {
		filter["created_at"] = bson.M{"$gte": filterOptions.StartDate, "$lte": filterOptions.EndDate}
	} else if filterOptions.StartDate != nil {
		filter["created_at"] = bson.M{"$gte": filterOptions.StartDate}
	} else if filterOptions.EndDate != nil {
		filter["created_at"] = bson.M{"$lte": filterOptions.EndDate}
	}

	// --- MOVE THE POPULARITY LOGIC HERE ---
	if filterOptions.MinPopularity != nil {
		filter["likes"] = bson.M{"$gte": *filterOptions.MinPopularity}
	}
	if filterOptions.MaxPopularity != nil {
		if _, ok := filter["likes"]; ok {
			filter["likes"].(bson.M)["$lte"] = *filterOptions.MaxPopularity
		} else {
			filter["likes"] = bson.M{"$lte": *filterOptions.MaxPopularity}
		}
	}
	// --- End of filter building ---

	// --- STAGE 2: Build the Query Options (Pagination & Sorting) ---
	findOptions := options.Find()
	findOptions.SetLimit(filterOptions.Limit)
	findOptions.SetSkip((filterOptions.Page - 1) * filterOptions.Limit)

	sort := bson.M{}
	switch filterOptions.SortBy {
	case "popularity":
		sort["likes"] = -1
	case "date_asc":
		sort["created_at"] = 1
	case "date_desc":
		sort["created_at"] = -1
	default:
		sort["created_at"] = -1
	}
	findOptions.SetSort(sort)

	// --- STAGE 3: Execute the Database Queries ---
	// The `filter` variable now contains ALL the required criteria.
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []entities.Blog
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}

	// The `filter` for counting is also now correct.
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// --- STAGE 4: Return the Results ---
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

// new

// FindByBlogAndUser finds a specific interaction record for a given blog and user.
func (r *mongoBlogInteractionRepository) FindByBlogAndUser(ctx context.Context, blogID, userID primitive.ObjectID) (*entities.BlogInteraction, error) {
	var interaction entities.BlogInteraction
	filter := bson.M{
		"blog_id": blogID,
		"user_id": userID,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&interaction)
	if err != nil {
		// It's normal for a document not to be found, so we don't treat mongo.ErrNoDocuments as a critical error here.
		if err == mongo.ErrNoDocuments {
			return nil, nil // No interaction found, return nil without an error.
		}
		return nil, err
	}
	return &interaction, nil
}

// UpdateCounts directly sets the like and dislike counts on a blog post.
func (r *mongoBlogRepository) UpdateCounts(ctx context.Context, blogID primitive.ObjectID, likes, dislikes int64) error {
	filter := bson.M{"_id": blogID}
	update := bson.M{
		"$set": bson.M{
			"likes":    likes,
			"dislikes": dislikes,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("blog post not found to update counts")
	}
	return nil
}

// Comment
func (r *mongoCommentRepository) Create(ctx context.Context, comment *entities.Comment) (*entities.Comment, error) {
	result, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}
	comment.ID = result.InsertedID.(primitive.ObjectID)
	return comment, nil
}

func (r *mongoBlogRepository) IncrementCommentCount(ctx context.Context, blogID primitive.ObjectID) error {
	filter := bson.M{"_id": blogID}
	// Use the $inc operator for an atomic and fast increment operation.
	update := bson.M{"$inc": bson.M{"comment_count": 1}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
