package handler

import (
	"net/http"
	"strings"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/gin-gonic/gin"
)

func UserLogoutHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
			return
		}
		
		// Extract token
		token := strings.TrimPrefix(auth, "Bearer ")
		token = strings.TrimSpace(token)
		
		l := logic.NewUserLogoutLogic(c.Request.Context(), svcCtx)
		resp, err := l.UserLogout(&types.UserLogoutRequest{}, token)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

