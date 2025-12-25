package middleware

import (
	"net/http"

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
	uc, err := helper.AnalyzeToken(auth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set("UserId", uc.Id)
	c.Set("UserIdentity", uc.Identity)
	c.Set("UserName", uc.Name)
	c.Next()
}
