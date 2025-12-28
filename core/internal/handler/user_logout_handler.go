package handler

import (
	"net/http"
	"strings"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"

	"github.com/gin-gonic/gin"
)

func UserLogoutHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get refresh token from request body
		var req types.UserLogoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// If no body, try to get from header (backward compatibility)
			auth := c.GetHeader("Authorization")
			if auth != "" {
				// Extract token from Authorization header (for backward compatibility)
				token := strings.TrimPrefix(auth, "Bearer ")
				token = strings.TrimSpace(token)
				req.RefreshToken = token
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
				return
			}
		}

		if req.RefreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
			return
		}

		l := logic.NewUserLogoutLogic(c.Request.Context(), svcCtx)
		resp, err := l.UserLogout(&req, req.RefreshToken)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
