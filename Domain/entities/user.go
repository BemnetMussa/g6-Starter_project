package entities

import "time"

// User represents a user document in MongoDB.
type User struct {
    ID           string       `bson:"_id,omitempty" json:"id"`
    FullName     string       `bson:"full_name" json:"full_name"`
    Username     string       `bson:"username" json:"username"`
    Email        string       `bson:"email" json:"email"`
    Password     string       `bson:"password" json:"-"` // bcrypt hashed, omit in JSON
    Role         string       `bson:"role" json:"role"` // "admin", "user"
    ProfileImage *string      `bson:"profile_image,omitempty" json:"profile_image,omitempty"`
    Bio          *string      `bson:"bio,omitempty" json:"bio,omitempty"`
    ContactInfo  *ContactInfo `bson:"contact_info,omitempty" json:"contact_info,omitempty"`
    CreatedAt    time.Time    `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time    `bson:"updated_at" json:"updated_at"`
}

type ContactInfo struct {
    Phone   *string `bson:"phone,omitempty" json:"phone,omitempty"`
    Address *string `bson:"address,omitempty" json:"address,omitempty"`
}