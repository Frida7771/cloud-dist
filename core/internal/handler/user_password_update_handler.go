package handler

import (
	"net/http"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"

	"github.com/gin-gonic/gin"
)

func UserPasswordUpdateHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserPasswordUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l := logic.NewUserPasswordUpdateLogic(c.Request.Context(), svcCtx)
		resp, err := l.UserPasswordUpdate(&req, c.GetString("UserIdentity"))
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
