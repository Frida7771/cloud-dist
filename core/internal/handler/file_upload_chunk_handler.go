package handler

import (
	"net/http"

	"cloud-dist/core/helper"
	"cloud-dist/core/internal/logic"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/svc"

	"github.com/gin-gonic/gin"
)

func FileUploadChunkHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.PostForm("key") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "key is empty"})
			return
		}
		if c.PostForm("upload_id") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "upload_id is empty"})
			return
		}
		if c.PostForm("part_number") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "part_number is empty"})
			return
		}

		etag, err := helper.S3PartUpload(c.Request)
		if err != nil {
			respondError(c, err)
			return
		}

		l := logic.NewFileUploadChunkLogic(c.Request.Context(), svcCtx)
		_, err = l.FileUploadChunk(&types.FileUploadChunkRequest{})
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, &types.FileUploadChunkReply{Etag: etag})
	}
}
