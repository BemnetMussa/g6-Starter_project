package usecases

import (
	"errors"
	// "time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/services"
	"g6_starter_project/Infrastructure/utils"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type userInterface interface {
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

// Login checks credentials and returns user + tokens if valid
func (u *UserUsecase) Login(user *entities.User) (*entities.User, *entities.Token, error) {
	if !utils.IsValidEmail(user.Email) {
		return nil, nil, errors.New("invalid email format")
	}

	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, nil, err
	}

	// Check if user is verified
	if !existingUser.IsVerified {
		return nil, nil, errors.New("account not verified. Please check your email and verify your account")
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