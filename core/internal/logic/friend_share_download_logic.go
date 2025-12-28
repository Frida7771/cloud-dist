package logic

import (
	"context"
	"errors"
	"io"
	"log"

	"cloud-dist/core/helper"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
	"cloud-dist/core/svc"

	"gorm.io/gorm"
)

type FriendShareDownloadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendShareDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendShareDownloadLogic {
	return &FriendShareDownloadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendShareDownloadLogic) FriendShareDownload(req *types.FriendShareDownloadRequest, userIdentity string) (io.ReadCloser, string, string, error) {
	// First, verify the friend share record exists
	fs := new(models.FriendShare)
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", req.ShareIdentity).
		First(fs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("share record not found")
		}
		log.Printf("[FriendShareDownload] Failed to query friend share: %v", err)
		return nil, "", "", err
	}

	// Verify that the user is either the sender or receiver
	if fs.FromUserIdentity != userIdentity && fs.ToUserIdentity != userIdentity {
		log.Printf("[FriendShareDownload] Access denied: User %s is not authorized to download share %s", userIdentity, req.ShareIdentity)
		return nil, "", "", errors.New("access denied: you are not authorized to download this file")
	}

	// Verify that both users are friends
	// Check if friendship exists (bidirectional check)
	var friendCount int64
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.Friend{}).
		Where("((user_identity = ? AND friend_identity = ?) OR (user_identity = ? AND friend_identity = ?)) AND status = ?",
			fs.FromUserIdentity, fs.ToUserIdentity,
			fs.ToUserIdentity, fs.FromUserIdentity,
			"active").
		Count(&friendCount).Error
	if err != nil {
		log.Printf("[FriendShareDownload] Failed to check friendship: %v", err)
		return nil, "", "", err
	}

	if friendCount == 0 {
		log.Printf("[FriendShareDownload] Access denied: Users %s and %s are not friends", fs.FromUserIdentity, fs.ToUserIdentity)
		return nil, "", "", errors.New("access denied: friendship relationship not found")
	}

	// Get file info from repository_pool
	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", fs.RepositoryIdentity).
		First(rp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("file not found")
		}
		log.Printf("[FriendShareDownload] Failed to query repository pool: %v", err)
		return nil, "", "", err
	}

	// Extract S3 key from path
	s3Key := rp.Path
	if s3Key == "" {
		return nil, "", "", errors.New("file path is empty")
	}

	// Check if path is a URL (old data format)
	if len(s3Key) > 7 && (s3Key[:7] == "http://" || (len(s3Key) > 8 && s3Key[:8] == "https://")) {
		log.Printf("[FriendShareDownload] Warning: Path is a URL, not a key. This file needs to be re-uploaded.")
		return nil, "", "", errors.New("file path format is outdated, please re-upload the file")
	}

	// Download file from S3
	fileData, headResp, err := helper.S3Download(s3Key)
	if err != nil {
		log.Printf("[FriendShareDownload] Failed to download from S3: %v, key=%s", err, s3Key)
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

	log.Printf("[FriendShareDownload] Successfully prepared file download: key=%s, name=%s, type=%s", s3Key, fileName, contentType)
	return fileData, fileName, contentType, nil
}
