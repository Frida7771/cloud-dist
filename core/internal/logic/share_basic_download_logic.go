package logic

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"cloud-disk/core/helper"
	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type ShareBasicDownloadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDownloadLogic {
	return &ShareBasicDownloadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDownloadLogic) ShareBasicDownload(req *types.ShareBasicDownloadRequest) (io.ReadCloser, string, string, error) {
	// First, verify the share link is valid
	sb := new(models.ShareBasic)
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", req.ShareIdentity).
		First(sb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("share link not found or invalid")
		}
		log.Printf("[ShareBasicDownload] Failed to query share basic: %v", err)
		return nil, "", "", err
	}

	// Check if share link has expired
	if sb.ExpiredTime > 0 {
		// ExpiredTime is in seconds, check if it has passed
		createdAt := sb.CreatedAt
		expiredAt := createdAt.Add(time.Duration(sb.ExpiredTime) * time.Second)
		if time.Now().After(expiredAt) {
			return nil, "", "", errors.New("share link has expired")
		}
	}

	// Get file info from repository_pool
	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", sb.RepositoryIdentity).
		First(rp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("file not found")
		}
		log.Printf("[ShareBasicDownload] Failed to query repository pool: %v", err)
		return nil, "", "", err
	}

	// Extract S3 key from path
	s3Key := rp.Path
	if s3Key == "" {
		return nil, "", "", errors.New("file path is empty")
	}

	// Check if path is a URL (old data format)
	if len(s3Key) > 7 && (s3Key[:7] == "http://" || (len(s3Key) > 8 && s3Key[:8] == "https://")) {
		log.Printf("[ShareBasicDownload] Warning: Path is a URL, not a key. This file needs to be re-uploaded.")
		return nil, "", "", errors.New("file path format is outdated, please re-upload the file")
	}

	// Download file from S3
	fileData, headResp, err := helper.S3Download(s3Key)
	if err != nil {
		log.Printf("[ShareBasicDownload] Failed to download from S3: %v, key=%s", err, s3Key)
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

	log.Printf("[ShareBasicDownload] Successfully prepared file download: key=%s, name=%s, type=%s", s3Key, fileName, contentType)
	return fileData, fileName, contentType, nil
}

