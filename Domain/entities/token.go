package entities

import "time"

type Token struct {
    ID           string    `json:"id" bson:"_id,omitempty"`
    UserID       string    `json:"user_id" bson:"user_id"`
    AccessToken  string    `json:"access_token" bson:"access_token"`
    RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at" bson:"expires_at"` // refresh token expiry
    CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}
