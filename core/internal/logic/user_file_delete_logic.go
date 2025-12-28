package logic

import (
	"context"
	"errors"
	"log"

	"cloud-dist/core/helper"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"

	"gorm.io/gorm"
)

type UserFileDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileDeleteLogic {
	return &UserFileDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteRequest, userIdentity string) (resp *types.UserFileDeleteReply, err error) {
	// Get user_repository record to find repository_identity
	ur := new(models.UserRepository)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).
		First(ur).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	// If it's a folder (no repository_identity), just delete user_repository record
	if ur.RepositoryIdentity == "" {
		err = l.svcCtx.DB.WithContext(l.ctx).
			Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).
			Delete(&models.UserRepository{}).Error
		return
	}

	// Get file info from repository_pool
	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", ur.RepositoryIdentity).
		First(rp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Repository pool record not found, just delete user_repository
			log.Printf("[UserFileDelete] Repository pool record not found for identity: %s", ur.RepositoryIdentity)
			err = l.svcCtx.DB.WithContext(l.ctx).
				Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).
				Delete(&models.UserRepository{}).Error
			return
		}
		return nil, err
	}

	// Update user storage capacity (subtract file size)
	if rp.Size > 0 {
		if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
			Where("identity = ?", userIdentity).
			UpdateColumn("now_volume", gorm.Expr("now_volume - ?", rp.Size)).Error; err != nil {
			log.Printf("[UserFileDelete] Failed to update user capacity: %v", err)
			return nil, err
		}
	}

	// Delete user_repository record
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).
		Delete(&models.UserRepository{}).Error
	if err != nil {
		log.Printf("[UserFileDelete] Failed to delete user_repository: %v", err)
		return nil, err
	}

	// Delete repository_pool record
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", ur.RepositoryIdentity).
		Delete(&models.RepositoryPool{}).Error
	if err != nil {
		log.Printf("[UserFileDelete] Failed to delete repository_pool: %v", err)
		// Continue to delete S3 file even if database delete fails
	}

	// Delete file from S3
	s3Key := rp.Path
	if s3Key != "" {
		// Check if path is a URL (old data format)
		if len(s3Key) > 7 && (s3Key[:7] == "http://" || (len(s3Key) > 8 && s3Key[:8] == "https://")) {
			log.Printf("[UserFileDelete] Warning: Path is a URL, not a key. Cannot delete from S3. key=%s", s3Key)
		} else {
			// Delete from S3
			if err = helper.S3Delete(s3Key); err != nil {
				log.Printf("[UserFileDelete] Failed to delete file from S3: %v, key=%s", err, s3Key)
				// Don't return error, just log it (file is already removed from database)
			} else {
				log.Printf("[UserFileDelete] Successfully deleted file from S3: key=%s", s3Key)
			}
		}
	}

	return
}
