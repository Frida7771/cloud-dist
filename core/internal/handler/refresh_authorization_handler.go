package handler

import (
	"net/http"
	"strings"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"

	"github.com/gin-gonic/gin"
)

func RefreshAuthorizationHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.RefreshAuthorizationRequest
		// Try to bind JSON first, but don't fail if body is empty
		c.ShouldBindJSON(&req)

		// Get refresh token from Authorization header (preferred method)
		auth := c.GetHeader("Authorization")
		var refreshToken string
		if auth != "" {
			refreshToken = strings.TrimPrefix(auth, "Bearer ")
			refreshToken = strings.TrimSpace(refreshToken)
		}

		// If not in header, try request body
		if refreshToken == "" {
			refreshToken = c.PostForm("refresh_token")
		}

		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
			return
		}

		l := logic.NewRefreshAuthorizationLogic(c.Request.Context(), svcCtx)
		resp, err := l.RefreshAuthorization(&req, refreshToken)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
