package handler

import (
	"net/http"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/svc"

	"github.com/gin-gonic/gin"
)

func ShareBasicSaveHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify user is authenticated
		userIdentity := c.GetString("UserIdentity")
		if userIdentity == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: please login first"})
			return
		}

		var req types.ShareBasicSaveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l := logic.NewShareBasicSaveLogic(c.Request.Context(), svcCtx)
		resp, err := l.ShareBasicSave(&req, userIdentity)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
