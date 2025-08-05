package usecases

import (
	"errors"
	"fmt"
	"time"

	"g6_starter_project/Domain/entities"
)

// UserManagementUsecase handles user role changes and admin-only user operations
type UserManagementUsecase struct {
	userRepo entities.UserRepository
}

// NewUserManagementUsecase initializes the user management usecase
func NewUserManagementUsecase(userRepo entities.UserRepository) *UserManagementUsecase {
	return &UserManagementUsecase{
		userRepo: userRepo,
	}
}

// PromoteUser upgrades a user to admin role
func (u *UserManagementUsecase) PromoteUser(userID string) (*entities.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if user.Role == "admin" {
		return nil, errors.New("user is already an admin")
	}

	user.Role = "admin"
	user.UpdatedAt = time.Now()

	updatedUser, err := u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to promote user: %v", err)
	}

	updatedUser.Password = ""
	return updatedUser, nil
}

// DemoteUser downgrades an admin to regular user
func (u *UserManagementUsecase) DemoteUser(userID string) (*entities.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if user.Role == "user" {
		return nil, errors.New("user is already a regular user")
	}

	user.Role = "user"
	user.UpdatedAt = time.Now()

	updatedUser, err := u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to demote user: %v", err)
	}

	updatedUser.Password = ""
	return updatedUser, nil
}

// GetAllUsers returns all registered users (admin-only)
func (u *UserManagementUsecase) GetAllUsers() ([]*entities.User, error) {
	return nil, errors.New("get all users not implemented yet")
}

// GetUserByID returns a user by ID
func (u *UserManagementUsecase) GetUserByID(userID string) (*entities.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}
