package usecases

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"g6-Starter_project/Domain/entities"
	"g6-Starter_project/Infrastructure/utils"
	"g6-Starter_project/Infrastructure/services"

)

type userInterface interface {
	Register(user *entities.User) (*entities.User, error)
	Login(user *entities.User) (*entities.User, error)
}

type UserUsecase struct {
	userRepo entities.UserRepository
}

func NewUserUsecase(userRepo entities.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) Register(user *entities.User) (*entities.User, error) {
	// Validate email format
	if !utils.IsValidEmail(user.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password requirements
	if !utils.IsValidPassword(user.Password) {
		return nil, errors.New("invalid password format - password must be at least 8 characters and contain uppercase, lowercase, number and special character")
	}

	// Check if email already exists
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Generate ObjectID for new user
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Set role based on whether this is the first user
	userCount, err := u.userRepo.GetUserCount()
	if err != nil {
		return nil, err
	}

	if userCount == 0 {
		// First user becomes admin
		user.Role = "admin"
	} else {
		// All other users are regular users
		user.Role = "user"
	}

	// Hash the password
	bcryptService := services.NewBcryptService(10)
	hashedPassword, err := bcryptService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// Create the user
	createdUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	
	// Clear password from response for security
	createdUser.Password = ""
	return createdUser, nil
}

func (u *UserUsecase) Login(user *entities.User) (*entities.User, error) {
	// Validate email format
	if !utils.IsValidEmail(user.Email) {
		return nil, errors.New("invalid email format")
	}
	// validate password format
	if !utils.IsValidPassword(user.Password) {
		return nil, errors.New("invalid password format - password must be at least 8 characters and contain uppercase, lowercase, number and special character")
	}

	// Check if user exists
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	
	// Check if password is correct
	bcryptService := services.NewBcryptService(10)
	err = bcryptService.ComparePassword(existingUser.Password, user.Password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// Clear password from response for security
	existingUser.Password = ""
	return existingUser, nil	


}