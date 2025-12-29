package logic

import (
	"context"
	"errors"
	"time"

	"cloud-dist/core/helper"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
	"cloud-dist/core/svc"
)

type ShareBasicDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDetailLogic {
	return &ShareBasicDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDetailLogic) ShareBasicDetail(req *types.ShareBasicDetailRequest) (resp *types.ShareBasicDetailReply, err error) {
	// Verify share link exists and not expired
	sb := new(models.ShareBasic)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", req.Identity).
		First(sb).Error
	if err != nil {
		return
	}

	// Check if share link has expired
	if sb.ExpiredTime > 0 {
		createdAt := sb.CreatedAt
		expiredAt := createdAt.Add(time.Duration(sb.ExpiredTime) * time.Second)
		if time.Now().After(expiredAt) {
			return nil, errors.New("share link has expired")
		}
	}

	resp = new(types.ShareBasicDetailReply)
	err = l.svcCtx.DB.WithContext(l.ctx).Table("share_basic").
		Select("share_basic.repository_identity, user_repository.name, repository_pool.ext, repository_pool.size, repository_pool.path").
		Joins("LEFT JOIN repository_pool ON share_basic.repository_identity = repository_pool.identity").
		Joins("LEFT JOIN user_repository ON user_repository.identity = share_basic.user_repository_identity").
		Where("share_basic.identity = ?", req.Identity).
		Take(resp).Error
	if err != nil {
		return
	}

	// Generate presigned URLs for preview and download
	// Presigned URL expiration should match share link expiration (3 days = 72 hours)
	// This allows preview and download without going through backend
	if resp.Path != "" {
		// Path is S3 key
		s3Key := resp.Path
		// Generate presigned URL for preview (no Content-Disposition)
		resp.Path = helper.S3PresignedURL(s3Key, 72)
		// Generate presigned URL for download (with Content-Disposition: attachment)
		fileName := resp.Name + resp.Ext
		resp.DownloadUrl = helper.S3PresignedURLWithDisposition(s3Key, 72, fileName)
	}

	return
}
