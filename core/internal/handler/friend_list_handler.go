package handler

import (
	"net/http"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"

	"github.com/gin-gonic/gin"
)

func FriendListHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.FriendListRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l := logic.NewFriendListLogic(c.Request.Context(), svcCtx)
		resp, err := l.FriendList(&req, c.GetString("UserIdentity"))
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
