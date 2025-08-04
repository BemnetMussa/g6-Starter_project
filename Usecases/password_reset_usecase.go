package usecases

import (
	"errors"
	"fmt"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/services"
	"g6_starter_project/Infrastructure/utils"
)

type PasswordResetUsecase struct {
	userRepo     entities.UserRepository
	jwtService   *services.JWTService
	emailService *services.EmailService
	rateLimiter  *services.RateLimiter
}

func NewPasswordResetUsecase(userRepo entities.UserRepository, jwtService *services.JWTService, emailService *services.EmailService, rateLimiter *services.RateLimiter) *PasswordResetUsecase {
	return &PasswordResetUsecase{
		userRepo:     userRepo,
		jwtService:   jwtService,
		emailService: emailService,
		rateLimiter:  rateLimiter,
	}
}

func (p *PasswordResetUsecase) RequestPasswordReset(email string) error {
	if !utils.IsValidEmail(email) {
		return errors.New("invalid email format")
	}

	// rate limit 3 request per email6
	if !p.rateLimiter.IsAllowed("forgot_password:"+email, 3, time.Hour) {
		return errors.New("too many password reset requests. please wait before trying again")
	}

	user, err := p.userRepo.GetUserByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	resetToken, err := p.jwtService.GenerateResetToken(user.ID.Hex())
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %v", err)
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	if err := p.userRepo.UpdateResetToken(user.ID.Hex(), &resetToken, &expiresAt); err != nil {
		return fmt.Errorf("failed to store reset token: %v", err)
	}

	if err := p.emailService.SendPasswordResetEmail(user.Email, user.FullName, resetToken); err != nil {
		return fmt.Errorf("failed to send reset email: %v", err)
	}

	return nil
}

// Reset user password 
func (p *PasswordResetUsecase) ResetPassword(token, newPassword string) error {
	if !utils.IsValidPassword(newPassword) {
		return errors.New("invalid password format - password must be at least 8 characters and contain uppercase, lowercase, number and special character")
	}

	if _, err := p.jwtService.ValidateResetToken(token); err != nil {
		return fmt.Errorf("invalid or expired reset token: %v", err)
	}

	user, err := p.userRepo.GetUserByResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid reset token: %v", err)
	}

	if user.ResetTokenExpiresAt != nil && time.Now().After(*user.ResetTokenExpiresAt) {
		return errors.New("reset token has expired")
	}

	bcryptService := services.NewBcryptService(10)
	hashedPassword, err := bcryptService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = hashedPassword
	user.ResetToken = nil
	user.ResetTokenExpiresAt = nil
	user.UpdatedAt = time.Now()

	if _, err := p.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	if err := p.emailService.SendPasswordChangeNotification(user.Email, user.FullName); err != nil {
		fmt.Printf("Warning: Failed to send password change notification: %v\n", err)
	}

	return nil
}
