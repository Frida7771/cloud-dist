package handler

import (
	"errors"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/svc"

	"github.com/gin-gonic/gin"
)

func FriendShareSaveHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.FriendShareSaveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondError(c, err)
			return
		}

		userIdentity := c.GetString("UserIdentity")
		if userIdentity == "" {
			respondUnauthorized(c, errors.New("unauthorized"))
			return
		}

		l := logic.NewFriendShareSaveLogic(c.Request.Context(), svcCtx)
		resp, err := l.FriendShareSave(&req, userIdentity)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(200, resp)
	}
}
