package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"g6-Starter_project/Domain/entities"
)

type UserRepositoryImpl struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Collection) entities.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *entities.User) (*entities.User, error) {
	_, err := r.db.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*entities.User, error) {
	filter := bson.M{"email": strings.ToLower(email)}
	var user entities.User
	err := r.db.FindOne(context.TODO(), filter).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserCount() (int64, error) {
	count, err := r.db.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepositoryImpl) GetUserByID(id string) (*entities.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	
	filter := bson.M{"_id": objectID}
	var user entities.User
	err = r.db.FindOne(context.TODO(), filter).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(user *entities.User) (*entities.User, error) {
	user.UpdatedAt = time.Now()
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	
	_, err := r.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) DeleteUser(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}
	
	filter := bson.M{"_id": objectID}
	_, err = r.db.DeleteOne(context.TODO(), filter)
	return err
}

