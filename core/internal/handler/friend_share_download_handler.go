package handler

import (
	"errors"
	"io"
	"log"
	"net/http"

	"cloud-dist/core/internal/logic"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/svc"

	"github.com/gin-gonic/gin"
)

func FriendShareDownloadHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[FriendShareDownloadHandler] Download request received")

		// Check authentication
		userIdentity := c.GetString("UserIdentity")
		if userIdentity == "" {
			log.Printf("[FriendShareDownloadHandler] No UserIdentity found in context")
			respondUnauthorized(c, errors.New("unauthorized: please login first"))
			return
		}
		log.Printf("[FriendShareDownloadHandler] User authenticated: %s", userIdentity)

		shareIdentity := c.Query("identity")
		if shareIdentity == "" {
			log.Printf("[FriendShareDownloadHandler] Missing share identity parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "share identity is required"})
			return
		}
		log.Printf("[FriendShareDownloadHandler] Requesting download for share: %s", shareIdentity)

		req := &types.FriendShareDownloadRequest{
			ShareIdentity: shareIdentity,
		}

		l := logic.NewFriendShareDownloadLogic(c.Request.Context(), svcCtx)
		fileData, fileName, contentType, err := l.FriendShareDownload(req, userIdentity)
		if err != nil {
			log.Printf("[FriendShareDownloadHandler] Failed to download file: %v", err)
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
			log.Printf("[FriendShareDownloadHandler] Failed to stream file: %v", err)
			return
		}
	}
}
