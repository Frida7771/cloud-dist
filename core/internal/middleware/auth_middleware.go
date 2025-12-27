package middleware

import (
	"net/http"
	"strings"

	"cloud-disk/core/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AuthMiddleware struct {
	RDB *redis.Client
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// SetRedisClient sets the Redis client for token blacklist checking
func (m *AuthMiddleware) SetRedisClient(rdb *redis.Client) {
	m.RDB = rdb
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Remove "Bearer " prefix if present
	token := strings.TrimPrefix(auth, "Bearer ")
	token = strings.TrimSpace(token)

	// Check if token is blacklisted (if Redis is available)
	if m.RDB != nil {
		blacklistKey := "token:blacklist:" + token
		_, err := m.RDB.Get(c.Request.Context(), blacklistKey).Result()
		if err == nil {
			// Token found in blacklist
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			return
		}
	}

	uc, err := helper.AnalyzeToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Verify token contains required fields
	if uc.Identity == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: user identity is empty"})
		return
	}

	// Set user information in context
	c.Set("UserId", uc.Id)
	c.Set("UserIdentity", uc.Identity)
	c.Set("UserName", uc.Name)

	c.Next()
}
