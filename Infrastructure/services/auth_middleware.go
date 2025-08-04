package services

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)
	
type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "userRole"
)

// JWTServiceInterface should define ValidateToken(token string) (map[string]interface{}, error)
type NewJWTServiceInterface interface {
	ValidateToken(token string) (map[string]interface{}, error)
}

// AuthMiddleware verifies JWT access tokens and sets user ID and role in Gin context
func AuthMiddleware(authSvc JWTServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
			}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims, err := authSvc.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
			return
		}
		role, _ := claims["role"].(string)

		// Store in Gin context
		c.Set(string(userIDKey), sub)
		c.Set(string(roleKey), role)

		c.Next()
	}
}

// RoleAuthorization middleware checks if user role is allowed
func RoleAuthorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(string(roleKey))
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role format"})
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
	}
}

// GetUserID retrieves user ID from standard request context
func GetUserID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(userIDKey).(string)
	return id, ok
}

// GetUserRole retrieves user role from standard request context
func GetUserRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(roleKey).(string)
	return role, ok
}

// GinGetUserID retrieves user ID from Gin context
func GinGetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(string(userIDKey))
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GinGetUserRole retrieves user role from Gin context
func GinGetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get(string(roleKey))
	if !exists {
		return "", false
	}
	role, ok := userRole.(string)
	return role, ok
}
