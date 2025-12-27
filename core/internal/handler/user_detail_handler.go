package handler

import (
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"

	"github.com/gin-gonic/gin"
)

func UserDetailHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user identity from context (set by auth middleware)
		userIdentity := c.GetString("UserIdentity")
		if userIdentity == "" {
			// Debug: check what's in the context
			userId := c.GetInt64("UserId")
			userName := c.GetString("UserName")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User identity not found",
				"debug": gin.H{
					"user_id":       userId,
					"user_name":     userName,
					"user_identity": userIdentity,
				},
			})
			return
		}

		req := &types.UserDetailRequest{
			Identity: userIdentity,
		}

		l := logic.NewUserDetailLogic(c.Request.Context(), svcCtx)
		resp, err := l.UserDetail(req)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
