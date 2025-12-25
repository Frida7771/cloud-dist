package handler

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
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
		log.Printf("[FileUpload] 开始处理文件上传请求")

		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			log.Printf("[FileUpload] 获取文件失败: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[FileUpload] 文件信息: 文件名=%s, 大小=%d 字节", fileHeader.Filename, fileHeader.Size)

		userIdentity := c.GetString("UserIdentity")
		log.Printf("[FileUpload] 用户身份: %s", userIdentity)

		ub := new(models.UserBasic)
		err = svcCtx.DB.WithContext(c.Request.Context()).
			Select("now_volume", "total_volume").
			Where("identity = ?", userIdentity).First(ub).Error
		if err != nil {
			log.Printf("[FileUpload] 查询用户容量失败: %v", err)
			respondError(c, err)
			return
		}
		log.Printf("[FileUpload] 用户容量: 已使用=%d 字节, 总容量=%d 字节", ub.NowVolume, ub.TotalVolume)

		if fileHeader.Size+ub.NowVolume > ub.TotalVolume {
			log.Printf("[FileUpload] 容量不足: 文件大小=%d, 已使用=%d, 总容量=%d", fileHeader.Size, ub.NowVolume, ub.TotalVolume)
			respondError(c, errors.New("已超出当前容量"))
			return
		}

		log.Printf("[FileUpload] 开始读取文件内容并计算 MD5")
		b := make([]byte, fileHeader.Size)
		if _, err = file.Read(b); err != nil {
			log.Printf("[FileUpload] 读取文件失败: %v", err)
			respondError(c, err)
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(b))
		log.Printf("[FileUpload] 文件 MD5: %s", hash)

		rp := new(models.RepositoryPool)
		err = svcCtx.DB.WithContext(c.Request.Context()).Where("hash = ?", hash).First(rp).Error
		if err == nil {
			log.Printf("[FileUpload] 文件已存在（秒传）: identity=%s", rp.Identity)
			c.JSON(http.StatusOK, &types.FileUploadReply{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[FileUpload] 查询文件记录失败: %v", err)
			respondError(c, err)
			return
		}

		log.Printf("[FileUpload] 开始上传到 S3")
		filePath, err := helper.S3Upload(c.Request)
		if err != nil {
			log.Printf("[FileUpload] S3 上传失败: %v", err)
			respondError(c, err)
			return
		}
		log.Printf("[FileUpload] S3 上传成功: 路径=%s", filePath)

		req := types.FileUploadRequest{
			Name: fileHeader.Filename,
			Ext:  path.Ext(fileHeader.Filename),
			Size: fileHeader.Size,
			Hash: hash,
			Path: filePath,
		}

		log.Printf("[FileUpload] 保存文件记录到数据库")
		l := logic.NewFileUploadLogic(c.Request.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			log.Printf("[FileUpload] 保存文件记录失败: %v", err)
			respondError(c, err)
			return
		}
		log.Printf("[FileUpload] 文件上传完成: identity=%s", resp.Identity)

		// 自动保存到用户仓库并更新容量
		log.Printf("[FileUpload] 自动保存到用户仓库")
		saveReq := types.UserRepositorySaveRequest{
			ParentId:           0, // 根目录
			RepositoryIdentity: resp.Identity,
			Ext:                resp.Ext,
			Name:               resp.Name,
		}
		saveLogic := logic.NewUserRepositorySaveLogic(c.Request.Context(), svcCtx)
		_, err = saveLogic.UserRepositorySave(&saveReq, userIdentity)
		if err != nil {
			log.Printf("[FileUpload] 保存到用户仓库失败: %v", err)
			// 即使保存到用户仓库失败，也返回上传成功（文件已经在中心存储池了）
			// 用户稍后可以手动保存
		} else {
			log.Printf("[FileUpload] 已保存到用户仓库，容量已更新")
		}

		c.JSON(http.StatusOK, resp)
	}
}
