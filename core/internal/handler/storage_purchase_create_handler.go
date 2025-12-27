package handler

import (
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/svc"

	"github.com/gin-gonic/gin"
)

func StoragePurchaseCreateHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.StoragePurchaseCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIdentity := c.GetString("UserIdentity")
		if userIdentity == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User identity not found"})
			return
		}

		l := logic.NewStoragePurchaseCreateLogic(c.Request.Context(), svcCtx)
		resp, err := l.StoragePurchaseCreate(&req, userIdentity)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}


