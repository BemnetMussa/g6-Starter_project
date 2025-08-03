package usecases

import (
	"errors"
	"fmt"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/services"
	"g6_starter_project/Infrastructure/utils"
)

// PasswordResetUsecase handles password reset logic
type PasswordResetUsecase struct {
	userRepo     entities.UserRepository
	jwtService   *services.JWTService
	emailService *services.EmailService
	rateLimiter  *services.RateLimiter
}

// NewPasswordResetUsecase creates a new password reset usecase
func NewPasswordResetUsecase(userRepo entities.UserRepository, jwtService *services.JWTService, emailService *services.EmailService, rateLimiter *services.RateLimiter) *PasswordResetUsecase {
	return &PasswordResetUsecase{
		userRepo:     userRepo,
		jwtService:   jwtService,
		emailService: emailService,
		rateLimiter:  rateLimiter,
	}
}

// RequestPasswordReset handles the forgot password request
func (p *PasswordResetUsecase) RequestPasswordReset(email string) error {
	// Validate email format
	if !utils.IsValidEmail(email) {
		return errors.New("invalid email format")
	}

	// Rate limiting: max 3 requests per email per hour
	if !p.rateLimiter.IsAllowed("forgot_password:"+email, 3, time.Hour) {
		return errors.New("too many password reset requests. please wait before trying again")
	}

	// Check if user exists
	user, err := p.userRepo.GetUserByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate reset token
	resetToken, err := p.jwtService.GenerateResetToken(user.ID.Hex())
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %v", err)
	}

	// Set expiry time (15 minutes)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Store reset token in database
	err = p.userRepo.UpdateResetToken(user.ID.Hex(), &resetToken, &expiresAt)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %v", err)
	}

	// Send reset email
	err = p.emailService.SendPasswordResetEmail(user.Email, user.FullName, resetToken)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %v", err)
	}

	return nil
}

// ResetPassword handles the password reset with token
func (p *PasswordResetUsecase) ResetPassword(token, newPassword string) error {
	// Validate new password
	if !utils.IsValidPassword(newPassword) {
		return errors.New("invalid password format - password must be at least 8 characters and contain uppercase, lowercase, number and special character")
	}

	// Validate reset token
	_, err := p.jwtService.ValidateResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token: %v", err)
	}

	// Get user from database
	user, err := p.userRepo.GetUserByResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid reset token: %v", err)
	}

	// Check if token is expired in database
	if user.ResetTokenExpiresAt != nil && time.Now().After(*user.ResetTokenExpiresAt) {
		return errors.New("reset token has expired")
	}

	// Hash new password
	bcryptService := services.NewBcryptService(10)
	hashedPassword, err := bcryptService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update user password and clear reset token
	user.Password = hashedPassword
	user.ResetToken = nil
	user.ResetTokenExpiresAt = nil
	user.UpdatedAt = time.Now()

	// Save updated user
	_, err = p.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	// Send password change notification email
	err = p.emailService.SendPasswordChangeNotification(user.Email, user.FullName)
	if err != nil {
		// Log error but don't fail the password reset
		fmt.Printf("Warning: Failed to send password change notification: %v\n", err)
	}

	return nil
} 