package usecases

import (
	"errors"
	"time"

	"g6_starter_project/Domain/entities"
	"g6_starter_project/Infrastructure/services"
)

type VerificationUsecase struct {
	userRepo    entities.UserRepository
	emailService *services.EmailService
}

func NewVerificationUsecase(userRepo entities.UserRepository, emailService *services.EmailService) *VerificationUsecase {
	return &VerificationUsecase{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

// RegisterWithVerification registers a user and sends verification email
func (v *VerificationUsecase) RegisterWithVerification(user *entities.User) (*entities.User, error) {
	// Set user as unverified initially
	user.IsVerified = false
	
	// Hash the password before storing
	bcryptService := services.NewBcryptService(10)
	hashedPassword, err := bcryptService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	
	// Generate verification token
	verificationToken, err := v.emailService.GenerateVerificationToken()
	if err != nil {
		return nil, err
	}
	
	// Store verification token in user (we'll use ResetToken field for this)
	user.ResetToken = &verificationToken
	
	// Set token expiration (24 hours)
	expiresAt := time.Now().Add(24 * time.Hour)
	user.ResetTokenExpiresAt = &expiresAt
	
	// Create user in database
	createdUser, err := v.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	
	// Send verification email
	username := user.UserName
	if username == "" {
		username = user.FullName
	}
	if username == "" {
		username = "User"
	}
	
	err = v.emailService.SendVerificationEmail(user.Email, username, verificationToken)
	if err != nil {
		// Log the error but don't fail registration
		// In production, you might want to handle this differently
		return nil, err
	}
	
	// Don't expose sensitive data
	createdUser.Password = ""
	createdUser.ResetToken = nil
	createdUser.ResetTokenExpiresAt = nil
	
	return createdUser, nil
}

// VerifyEmail verifies a user's email using the verification token
func (v *VerificationUsecase) VerifyEmail(token string) error {
	// Find user by verification token
	user, err := v.userRepo.GetUserByResetToken(token)
	if err != nil {
		return errors.New("invalid verification token")
	}
	
	// Check if token is expired
	if user.ResetTokenExpiresAt != nil && time.Now().After(*user.ResetTokenExpiresAt) {
		return errors.New("verification token has expired")
	}
	
	// Update user verification status
	err = v.userRepo.UpdateVerificationStatus(user.ID.Hex(), true)
	if err != nil {
		return err
	}
	
	// Clear the verification token
	err = v.userRepo.UpdateResetToken(user.ID.Hex(), nil, nil)
	if err != nil {
		return err
	}
	
	// Send welcome email
	username := user.UserName
	if username == "" {
		username = user.FullName
	}
	if username == "" {
		username = "User"
	}
	
	err = v.emailService.SendWelcomeEmail(user.Email, username)
	if err != nil {
		// Log the error but don't fail verification
		// In production, you might want to handle this differently
	}
	
	return nil
}

// ResendVerificationEmail resends verification email to user
func (v *VerificationUsecase) ResendVerificationEmail(email string) error {
	// Find user by email
	user, err := v.userRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	
	// Check if user is already verified
	if user.IsVerified {
		return errors.New("user is already verified")
	}
	
	// Generate new verification token
	verificationToken, err := v.emailService.GenerateVerificationToken()
	if err != nil {
		return err
	}
	
	// Set token expiration (24 hours)
	expiresAt := time.Now().Add(24 * time.Hour)
	
	// Update user with new verification token
	err = v.userRepo.UpdateResetToken(user.ID.Hex(), &verificationToken, &expiresAt)
	if err != nil {
		return err
	}
	
	// Send verification email
	username := user.UserName
	if username == "" {
		username = user.FullName
	}
	if username == "" {
		username = "User"
	}
	
	err = v.emailService.SendVerificationEmail(user.Email, username, verificationToken)
	if err != nil {
		return err
	}
	
	return nil
} 