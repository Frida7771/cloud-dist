package handler

import (
	"errors"
	"io"
	"log"
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/svc"

	"github.com/gin-gonic/gin"
)

func FileDownloadHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[FileDownloadHandler] Download request received")

		// Check authentication first
		userIdentity := c.GetString("UserIdentity")
		if userIdentity == "" {
			log.Printf("[FileDownloadHandler] No UserIdentity found in context - authentication may have failed")
			// Check if Authorization header is present
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				log.Printf("[FileDownloadHandler] No Authorization header found")
			} else {
				log.Printf("[FileDownloadHandler] Authorization header present but UserIdentity not set")
			}
			respondUnauthorized(c, errors.New("unauthorized: please login first"))
			return
		}
		log.Printf("[FileDownloadHandler] User authenticated: %s", userIdentity)

		repositoryIdentity := c.Query("identity")
		if repositoryIdentity == "" {
			log.Printf("[FileDownloadHandler] Missing repository identity parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "repository identity is required"})
			return
		}
		log.Printf("[FileDownloadHandler] Requesting download for repository: %s", repositoryIdentity)

		l := logic.NewFileDownloadLogic(c.Request.Context(), svcCtx)
		fileData, fileName, contentType, err := l.FileDownload(repositoryIdentity, userIdentity)
		if err != nil {
			log.Printf("[FileDownloadHandler] Failed to download file: %v", err)
			// Check if it's an access denied error
			if err.Error() == "access denied: file not found in your repository" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			} else {
				respondError(c, err)
			}
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
			log.Printf("[FileDownload] Failed to stream file: %v", err)
			return
		}
	}
}
