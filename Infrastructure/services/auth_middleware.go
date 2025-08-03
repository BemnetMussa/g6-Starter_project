package services

import (
	"context"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
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

// GinAuthMiddleware is the Gin-compatible version of AuthMiddleware
func GinAuthMiddleware(authSvc JWTServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate the access token using JWTService
		claims, err := authSvc.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user ID ("sub" claim) from token
		sub, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
			c.Abort()
			return
		}

		// Extract user role (may be empty)
		role, _ := claims["role"].(string)

		// Store user ID and role in Gin context
		c.Set("userID", sub)
		c.Set("userRole", role)

		c.Next()
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

// GinRoleAuthorization is the Gin-compatible version of RoleAuthorization
func GinRoleAuthorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
		c.Abort()
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

// GinGetUserID gets user ID from Gin context
func GinGetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GinGetUserRole gets user role from Gin context
func GinGetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("userRole")
	if !exists {
		return "", false
	}
	role, ok := userRole.(string)
	return role, ok
}
