package handler

import (
	"io"
	"log"
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/svc"

	"github.com/gin-gonic/gin"
)

func ShareBasicDownloadHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[ShareBasicDownloadHandler] Download request received")

		shareIdentity := c.Query("identity")
		if shareIdentity == "" {
			log.Printf("[ShareBasicDownloadHandler] Missing share identity parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "share identity is required"})
			return
		}
		log.Printf("[ShareBasicDownloadHandler] Requesting download for share: %s", shareIdentity)

		req := &types.ShareBasicDownloadRequest{
			ShareIdentity: shareIdentity,
		}

		l := logic.NewShareBasicDownloadLogic(c.Request.Context(), svcCtx)
		fileData, fileName, contentType, err := l.ShareBasicDownload(req)
		if err != nil {
			log.Printf("[ShareBasicDownloadHandler] Failed to download file: %v", err)
			respondError(c, err)
			return
		}
		defer fileData.Close()

		// Set headers for file download
		c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		c.Header("Content-Type", contentType)
		c.Header("Content-Transfer-Encoding", "binary")

		// Stream file content to response
		_, err = io.Copy(c.Writer, fileData)
		if err != nil {
			log.Printf("[ShareBasicDownloadHandler] Failed to stream file: %v", err)
			return
		}
	}
}

