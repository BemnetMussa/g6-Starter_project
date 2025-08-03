package usecases

import (
	"errors"
	"fmt"

	"g6_starter_project/Domain/entities"
)

// UserManagementUsecase handles admin user management operations
type UserManagementUsecase struct {
	userRepo entities.UserRepository
}

// NewUserManagementUsecase creates a new user management usecase
func NewUserManagementUsecase(userRepo entities.UserRepository) *UserManagementUsecase {
	return &UserManagementUsecase{
		userRepo: userRepo,
	}
}

// PromoteUser promotes a regular user to admin role
func (u *UserManagementUsecase) PromoteUser(userID string) (*entities.User, error) {
	// Get user by ID
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Check if user is already an admin
	if user.Role == "admin" {
		return nil, errors.New("user is already an admin")
	}

	// Promote user to admin
	user.Role = "admin"
	user.UpdatedAt = user.UpdatedAt

	// Update user in database
	updatedUser, err := u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to promote user: %v", err)
	}

	// Clear password from response for security
	updatedUser.Password = ""
	return updatedUser, nil
}

// DemoteUser demotes an admin user to regular user role
func (u *UserManagementUsecase) DemoteUser(userID string) (*entities.User, error) {
	// Get user by ID
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Check if user is already a regular user
	if user.Role == "user" {
		return nil, errors.New("user is already a regular user")
	}

	// Demote user to regular user
	user.Role = "user"
	user.UpdatedAt = user.UpdatedAt

	// Update user in database
	updatedUser, err := u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to demote user: %v", err)
	}

	// Clear password from response for security
	updatedUser.Password = ""
	return updatedUser, nil
}

// GetAllUsers returns all users (admin only)
func (u *UserManagementUsecase) GetAllUsers() ([]*entities.User, error) {
	// This would need to be implemented in the repository
	// For now, we'll return an error indicating this needs to be implemented
	return nil, errors.New("get all users not implemented yet")
}

// GetUserByID returns a specific user by ID
func (u *UserManagementUsecase) GetUserByID(userID string) (*entities.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Clear password from response for security
	user.Password = ""
	return user, nil
} 