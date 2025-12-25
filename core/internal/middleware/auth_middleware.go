package middleware

import (
	"net/http"
	"strings"

	"cloud-disk/core/helper"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
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
	uc, err := helper.AnalyzeToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set("UserId", uc.Id)
	c.Set("UserIdentity", uc.Identity)
	c.Set("UserName", uc.Name)
	c.Next()
}
