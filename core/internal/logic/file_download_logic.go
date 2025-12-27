package logic

import (
	"context"
	"errors"
	"io"
	"log"

	"cloud-disk/core/helper"
	"cloud-disk/core/svc"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type FileDownloadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadLogic) FileDownload(repositoryIdentity, userIdentity string) (io.ReadCloser, string, string, error) {
	// First, get file info from repository_pool to check if file exists
	rp := new(models.RepositoryPool)
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", repositoryIdentity).
		First(rp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("file not found")
		}
		log.Printf("[FileDownload] Failed to query repository pool: %v", err)
		return nil, "", "", err
	}

	// Check if user has access to this file via user_repository
	// Use GORM model for more reliable query
	var ur models.UserRepository
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_identity = ?", userIdentity).
		Where("repository_identity = ?", repositoryIdentity).
		First(&ur).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[FileDownload] Access denied: User %s does not have repository %s in user_repository", userIdentity, repositoryIdentity)

			// Additional debug info
			var totalCount int64
			l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
				Where("user_identity = ?", userIdentity).
				Count(&totalCount)
			log.Printf("[FileDownload] Debug: User %s has %d total files in repository", userIdentity, totalCount)

			var repoCount int64
			l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
				Where("repository_identity = ?", repositoryIdentity).
				Count(&repoCount)
			log.Printf("[FileDownload] Debug: Repository %s exists in %d user_repository records", repositoryIdentity, repoCount)

			return nil, "", "", errors.New("access denied: file not found in your repository")
		}
		log.Printf("[FileDownload] Failed to check access: %v", err)
		return nil, "", "", err
	}

	log.Printf("[FileDownload] Access granted: User %s has access to repository %s", userIdentity, repositoryIdentity)

	// Extract S3 key from path
	// Path should be an S3 key (e.g., cloud-disk/xxx.jpg)
	s3Key := rp.Path
	if s3Key == "" {
		return nil, "", "", errors.New("file path is empty")
	}

	// Check if path is a URL (old data format) - this should not happen with new uploads
	if len(s3Key) > 7 && (s3Key[:7] == "http://" || (len(s3Key) > 8 && s3Key[:8] == "https://")) {
		// This is an old URL format, we can't extract the key
		// User needs to re-upload the file
		log.Printf("[FileDownload] Warning: Path is a URL, not a key. This file needs to be re-uploaded.")
		return nil, "", "", errors.New("file path format is outdated, please re-upload the file")
	}

	// Download file from S3
	fileData, headResp, err := helper.S3Download(s3Key)
	if err != nil {
		log.Printf("[FileDownload] Failed to download from S3: %v, key=%s", err, s3Key)
		return nil, "", "", err
	}

	// Get file name and content type
	fileName := rp.Name
	if fileName == "" {
		fileName = "download"
		if rp.Ext != "" {
			fileName += rp.Ext
		}
	}

	contentType := "application/octet-stream"
	if headResp.ContentType != nil {
		contentType = *headResp.ContentType
	}

	log.Printf("[FileDownload] Successfully prepared file download: key=%s, name=%s, type=%s", s3Key, fileName, contentType)
	return fileData, fileName, contentType, nil
}
