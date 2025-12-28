package handler

import (
	"net/http"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"

	"github.com/gin-gonic/gin"
)

func ShareBasicDetailHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := types.ShareBasicDetailRequest{
			Identity: c.Query("identity"),
		}

		l := logic.NewShareBasicDetailLogic(c.Request.Context(), svcCtx)
		resp, err := l.ShareBasicDetail(&req)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
