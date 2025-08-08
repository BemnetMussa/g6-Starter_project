package repositories

import (
	"context"
	"fmt"

	"g6_starter_project/Domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenRepository struct {
	collection *mongo.Collection
}

func NewTokenRepository(col *mongo.Collection) *TokenRepository {
	return &TokenRepository{collection: col}
}


// Create inserts a new token into the database
func (r *TokenRepository) Create(ctx context.Context, token *entities.Token) error {
	result, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		return err
	}
	
	// Set the generated ID back to the token object
	// Convert ObjectID to string
	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		token.ID = objectID.Hex()
	}
	return nil
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
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
