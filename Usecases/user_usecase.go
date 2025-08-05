package usecases

import (
	"errors"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/services"
	"g6_starter_project/Infrastructure/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userInterface interface {
	Register(user *entities.User) (*entities.User, error)
	Login(user *entities.User) (*entities.User, error)
}

type UserUsecase struct {
	userRepo     entities.UserRepository
	tokenUsecase *TokenUsecase
}

func NewUserUsecase(userRepo entities.UserRepository, tokenUsecase *TokenUsecase) *UserUsecase {
	return &UserUsecase{
		userRepo:     userRepo,
		tokenUsecase: tokenUsecase,
	}
}

// Register validates user input, assigns role, hashes password, and stores the user
func (u *UserUsecase) Register(user *entities.User) (*entities.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return nil, errors.New("invalid email format")
	}

	if !utils.IsValidPassword(user.Password) {
		return nil, errors.New("invalid password format - password must be at least 8 characters and contain uppercase, lowercase, number and special character")
	}

	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Assign admin role to the very first user in the system
	userCount, err := u.userRepo.GetUserCount()
	if err != nil {
		return nil, err
	}
	if userCount == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	// Secure the password before storing it
	bcryptService := services.NewBcryptService(10)
	hashedPassword, err := bcryptService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	createdUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Don't expose password in the response
	createdUser.Password = ""
	return createdUser, nil
}

// Login checks credentials and returns user + tokens if valid
func (u *UserUsecase) Login(user *entities.User) (*entities.User, *entities.Token, error) {
	if !utils.IsValidEmail(user.Email) {
		return nil, nil, errors.New("invalid email format")
	}

	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, nil, err
	}

	// Compare entered password with stored hash
	bcryptService := services.NewBcryptService(10)
	err = bcryptService.ComparePassword(existingUser.Password, user.Password)
	if err != nil {
		return nil, nil, errors.New("invalid password")
	}

	// Generate JWT access & refresh tokens
	token, err := u.tokenUsecase.GenerateTokens(existingUser.ID.Hex(), existingUser.Role)
	if err != nil {
		return nil, nil, err
	}

	existingUser.Password = ""
	return existingUser, token, nil
}

// logout user
func (u *UserUsecase) Logout(userID string) error {
	return u.tokenUsecase.Logout(userID)
}
