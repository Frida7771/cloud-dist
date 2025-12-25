package handler

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"path"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIdentity := c.GetString("UserIdentity")
		ub := new(models.UserBasic)
		err = svcCtx.DB.WithContext(c.Request.Context()).
			Select("now_volume", "total_volume").
			Where("identity = ?", userIdentity).First(ub).Error
		if err != nil {
			respondError(c, err)
			return
		}
		if fileHeader.Size+ub.NowVolume > ub.TotalVolume {
			respondError(c, errors.New("已超出当前容量"))
			return
		}

		b := make([]byte, fileHeader.Size)
		if _, err = file.Read(b); err != nil {
			respondError(c, err)
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(b))

		rp := new(models.RepositoryPool)
		err = svcCtx.DB.WithContext(c.Request.Context()).Where("hash = ?", hash).First(rp).Error
		if err == nil {
			c.JSON(http.StatusOK, &types.FileUploadReply{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(c, err)
			return
		}

		filePath, err := helper.S3Upload(c.Request)
		if err != nil {
			respondError(c, err)
			return
		}

		req := types.FileUploadRequest{
			Name: fileHeader.Filename,
			Ext:  path.Ext(fileHeader.Filename),
			Size: fileHeader.Size,
			Hash: hash,
			Path: filePath,
		}

		l := logic.NewFileUploadLogic(c.Request.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
