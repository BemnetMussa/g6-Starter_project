package entities

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user document in MongoDB.
type User struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    FullName     string            `bson:"full_name" json:"full_name"`
    Username     string            `bson:"username" json:"username"`
    Email        string            `bson:"email" json:"email"`
    Password     string            `bson:"password" json:"password"` // Allow password during registration
    Role         string            `bson:"role" json:"role,omitempty"` // "admin", "user" - optional in JSON
    ProfileImage *string           `bson:"profile_image,omitempty" json:"profile_image,omitempty"`
    Bio          *string           `bson:"bio,omitempty" json:"bio,omitempty"`
    ContactInfo  *ContactInfo      `bson:"contact_info,omitempty" json:"contact_info,omitempty"`
    ResetToken   *string           `bson:"reset_token,omitempty" json:"reset_token,omitempty"`
    ResetTokenExpiresAt *time.Time `bson:"reset_token_expires_at,omitempty" json:"reset_token_expires_at,omitempty"`
    CreatedAt    time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time         `bson:"updated_at" json:"updated_at"`
}

type ContactInfo struct {
    Phone   *string `bson:"phone,omitempty" json:"phone,omitempty"`
    Address *string `bson:"address,omitempty" json:"address,omitempty"`
}

// interface for repository to use 
type UserRepository interface {
	CreateUser(user *User) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserCount() (int64, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(id string) error
	UpdateResetToken(userID string, resetToken *string, expiresAt *time.Time) error
	GetUserByResetToken(resetToken string) (*User, error)
}