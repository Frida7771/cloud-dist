package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/logic"
	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[FileUpload] Starting file upload request")

		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			log.Printf("[FileUpload] Failed to get file: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[FileUpload] File info: filename=%s, size=%d bytes", fileHeader.Filename, fileHeader.Size)

		userIdentity := c.GetString("UserIdentity")
		log.Printf("[FileUpload] User identity: %s", userIdentity)

		ub := new(models.UserBasic)
		err = svcCtx.DB.WithContext(c.Request.Context()).
			Select("now_volume", "total_volume").
			Where("identity = ?", userIdentity).First(ub).Error
		if err != nil {
			log.Printf("[FileUpload] Failed to query user capacity: %v", err)
			respondError(c, err)
			return
		}
		log.Printf("[FileUpload] User capacity: used=%d bytes, total=%d bytes", ub.NowVolume, ub.TotalVolume)

		if fileHeader.Size+ub.NowVolume > ub.TotalVolume {
			log.Printf("[FileUpload] Insufficient capacity: file size=%d, used=%d, total=%d", fileHeader.Size, ub.NowVolume, ub.TotalVolume)
			respondError(c, errors.New("storage capacity exceeded"))
			return
		}

		log.Printf("[FileUpload] Starting to read file content and calculate xxHash64")
		b := make([]byte, fileHeader.Size)
		if _, err = file.Read(b); err != nil {
			log.Printf("[FileUpload] Failed to read file: %v", err)
			respondError(c, err)
			return
		}
		hash := fmt.Sprintf("%016x", xxhash.Sum64(b))
		log.Printf("[FileUpload] File xxHash64: %s", hash)

		rp := new(models.RepositoryPool)
		err = svcCtx.DB.WithContext(c.Request.Context()).Where("hash = ?", hash).First(rp).Error
		if err == nil {
			log.Printf("[FileUpload] File already exists (instant upload): identity=%s", rp.Identity)
			c.JSON(http.StatusOK, &types.FileUploadReply{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[FileUpload] Failed to query file record: %v", err)
			respondError(c, err)
			return
		}

		log.Printf("[FileUpload] Starting upload to S3")
		filePath, err := helper.S3Upload(c.Request)
		if err != nil {
			log.Printf("[FileUpload] S3 upload failed: %v", err)
			respondError(c, err)
			return
		}
		log.Printf("[FileUpload] S3 upload successful: path=%s", filePath)

		req := types.FileUploadRequest{
			Name: fileHeader.Filename,
			Ext:  path.Ext(fileHeader.Filename),
			Size: fileHeader.Size,
			Hash: hash,
			Path: filePath,
		}

		log.Printf("[FileUpload] Saving file record to database")
		l := logic.NewFileUploadLogic(c.Request.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			log.Printf("[FileUpload] Failed to save file record: %v", err)
			respondError(c, err)
			return
		}
		log.Printf("[FileUpload] File upload completed: identity=%s", resp.Identity)

		// Don't automatically save to user repository
		// Frontend will call /user/repository/save with the selected folder
		// This allows users to choose which folder to save the file to

		c.JSON(http.StatusOK, resp)
	}
}
