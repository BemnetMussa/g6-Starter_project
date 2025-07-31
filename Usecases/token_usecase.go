package usecases

import (
	"context"
	"fmt"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/mongodb/repositories"
	"g6_starter_project/Infrastructure/services"

	"github.com/google/uuid" // For generating unique token ID
)

// TokenUsecase handles token logic between services and database
type TokenUsecase struct {
	repo       *repositories.TokenRepository
	jwtService *services.JWTService
}

// NewTokenUsecase creates a new usecase instance
func NewTokenUsecase(repo *repositories.TokenRepository, jwtService *services.JWTService) *TokenUsecase {
	return &TokenUsecase{
		repo:       repo,
		jwtService: jwtService,
	}
}

// GenerateTokens creates new access & refresh tokens for a user
func (u *TokenUsecase) GenerateTokens(userID, userRole string) (*entities.Token, error) {
	now := time.Now()

	// Create access and refresh tokens
	accessToken, refreshToken, err := u.jwtService.GenerateTokens(userID, userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	// Prepare token entity to store in DB
	token := &entities.Token{
		ID:           uuid.NewString(),
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(7 * 24 * time.Hour), // Token expiry (7 days)
		CreatedAt:    now,
	}

	err = u.repo.Create(context.Background(), token)
	if err != nil {
		return nil, fmt.Errorf("failed to store token: %v", err)
	}

	return token, nil
}

// RefreshToken validates existing refresh token and issues new tokens
func (u *TokenUsecase) RefreshToken(userID string) (*entities.Token, error) {
	// Get token from DB
	token, err := u.repo.FindByUserID(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find token: %v", err)
	}

	// Validate the refresh token (is it expired, tampered, etc.)
	_, err = u.jwtService.ValidateToken(token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	// Generate new access/refresh tokens and store again
	return u.GenerateTokens(userID, "user")
}
