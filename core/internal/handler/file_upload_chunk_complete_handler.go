package handler

import (
	"errors"
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"cloud-disk/core/svc"

	"github.com/gin-gonic/gin"
)

func FileUploadChunkCompleteHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.FileUploadChunkCompleteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIdentity := c.GetString("UserIdentity")
		ub := new(models.UserBasic)
		err := svcCtx.DB.WithContext(c.Request.Context()).
			Select("now_volume", "total_volume").
			Where("identity = ?", userIdentity).First(ub).Error
		if err != nil {
			respondError(c, err)
			return
		}
		if req.Size+ub.NowVolume > ub.TotalVolume {
			respondError(c, errors.New("已超出当前容量"))
			return
		}

		l := logic.NewFileUploadChunkCompleteLogic(c.Request.Context(), svcCtx)
		resp, err := l.FileUploadChunkComplete(&req)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
