package usecases

import (
	"fmt"
	"time"

	"g6_starter_project/Domain/entities"
)

// UserProfileUsecase handles user profile operations
type UserProfileUsecase struct {
	userRepo entities.UserRepository
}

// NewUserProfileUsecase creates a new user profile usecase
func NewUserProfileUsecase(userRepo entities.UserRepository) *UserProfileUsecase {
	return &UserProfileUsecase{
		userRepo: userRepo,
	}
}

// GetUserProfileByID gets a user profile by ID (internal use)
func (u *UserProfileUsecase) GetUserProfileByID(userID string) (*entities.User, error) {
	// Get user by ID
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Clear password from response for security
	user.Password = ""
	return user, nil
}

// UpdateUserProfile updates a user's profile information
func (u *UserProfileUsecase) UpdateUserProfile(userID string, updateData *entities.User) (*entities.User, error) {
	// Get existing user by ID
	existingUser, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Update only allowed fields (profile-related fields only)
	if updateData.Username != "" {
		existingUser.Username = updateData.Username
	}
	if updateData.FullName != "" {
		existingUser.FullName = updateData.FullName
	}
	if updateData.ProfileImage != nil {
		existingUser.ProfileImage = updateData.ProfileImage
	}
	if updateData.Bio != nil {
		existingUser.Bio = updateData.Bio
	}
	if updateData.ContactInfo != nil {
		existingUser.ContactInfo = updateData.ContactInfo
	}

	// Note: We intentionally do NOT update:
	// - Password (handled by password reset)
	// - Email (handled by separate email update system)
	// - Role (handled by admin promotion/demotion)

	// Update timestamp
	existingUser.UpdatedAt = time.Now()

	// Update user in database
	updatedUser, err := u.userRepo.UpdateUser(existingUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Clear password from response for security
	updatedUser.Password = ""
	return updatedUser, nil
}

