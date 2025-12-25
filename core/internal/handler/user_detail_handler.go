package handler

import (
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/gin-gonic/gin"
)

func UserDetailHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserDetailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l := logic.NewUserDetailLogic(c.Request.Context(), svcCtx)
		resp, err := l.UserDetail(&req)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
