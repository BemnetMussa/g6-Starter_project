package services

import (
	"context"
	"net/http"
	"strings"
        
)

// Define context keys to store user info in request context
type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "userRole"
)

// AuthMiddleware verifies JWT access tokens on incoming HTTP requests
func AuthMiddleware(authSvc JWTServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate the access token using JWTService
			claims, err := authSvc.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// // Optional: check if token is revoked by jti (token ID)
			// if jti, ok := claims["jti"].(string); ok {
			// 	if authSvc.IsTokenRevoked(jti) {
			// 		http.Error(w, "Token revoked", http.StatusUnauthorized)
			// 		return
			// 	}
			// }

			// Extract user ID ("sub" claim) from token
			sub, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid token subject", http.StatusUnauthorized)
				return
			}

			// Extract user role (may be empty)
			role, _ := claims["role"].(string)

			// Store user ID and role in request context
			ctx := context.WithValue(r.Context(), UserIDKey, sub)
			ctx = context.WithValue(ctx, RoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleAuthorization middleware checks if user role is allowed for the route
func RoleAuthorization(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetUserRole(r)
			if !ok {
				http.Error(w, "Role not found", http.StatusUnauthorized)
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

// Helper function to get user ID from request context
func GetUserID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(UserIDKey).(string)
	return id, ok
}

// Helper function to get user role from request context
func GetUserRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(RoleKey).(string)
	return role, ok
}
