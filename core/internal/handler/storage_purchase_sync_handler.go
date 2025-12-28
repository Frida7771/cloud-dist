package handler

import (
	"net/http"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/svc"

	"github.com/gin-gonic/gin"
)

func StoragePurchaseSyncHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.StoragePurchaseSyncRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.SessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
			return
		}

		l := logic.NewStoragePurchaseSyncLogic(c.Request.Context(), svcCtx)
		resp, err := l.StoragePurchaseSync(&req)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
