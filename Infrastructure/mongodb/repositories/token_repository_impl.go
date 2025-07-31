package repositories

import (
	"context"
	"fmt"

	"starter_project/Domain/entities"         
	"go.mongodb.org/mongo-driver/bson"        
	"go.mongodb.org/mongo-driver/mongo"       
)

type TokenRepository struct {
	collection *mongo.Collection 
}

// NewTokenRepository creates a new TokenRepository using the "tokens" collection
func NewTokenRepository(db *mongo.Database) *TokenRepository {
	return &TokenRepository{
		collection: db.Collection("tokens"),
	}
}

// Create inserts a new token into the database
func (r *TokenRepository) Create(ctx context.Context, token *entities.Token) error {
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

// FindByUserID retrieves a token by user ID
func (r *TokenRepository) FindByUserID(ctx context.Context, userID string) (*entities.Token, error) {
	var token entities.Token
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No token found for the given use
			return nil, fmt.Errorf("token not found for user %s", userID)
		}
		
		return nil, fmt.Errorf("failed to find token: %v", err)
	}
	return &token, nil
}

// Update modifies token fields for a specific user
func (r *TokenRepository) Update(ctx context.Context, userID string, update bson.M) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID}, 
		bson.M{"$set": update},    
	)
	return err
}

// DeleteByUserID removes a token by user ID
func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}
