package logic

import (
	"context"
	"errors"
	"log"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type UserRepositorySaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRepositorySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRepositorySaveLogic {
	return &UserRepositorySaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRepositorySaveLogic) UserRepositorySave(req *types.UserRepositorySaveRequest, userIdentity string) (resp *types.UserRepositorySaveReply, err error) {
	// Check if this file already exists in user_repository for this user (global deduplication - user level, not folder level)
	log.Printf("[UserRepositorySave] Checking for duplicate: user=%s, repository_identity=%s",
		userIdentity, req.RepositoryIdentity)

	existingUr := new(models.UserRepository)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_identity = ?", userIdentity).
		Where("repository_identity = ?", req.RepositoryIdentity).
		Where("deleted_at IS NULL").
		First(existingUr).Error
	if err == nil {
		// File already exists in user's repository (anywhere)
		log.Printf("[UserRepositorySave] File already exists in user repository: user=%s, repository_identity=%s, existing_identity=%s, existing_id=%d, existing_parent_id=%d",
			userIdentity, req.RepositoryIdentity, existingUr.Identity, existingUr.ID, existingUr.ParentId)
		// Return error to inform frontend that file already exists
		err = errors.New("file already exists")
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// File doesn't exist, proceed with save
		log.Printf("[UserRepositorySave] No duplicate found, proceeding with save")
		err = nil // Clear the error
	} else {
		// Database error
		log.Printf("[UserRepositorySave] Database error while checking for existing file: %v", err)
		return
	}

	// Check if file exceeds capacity
	rp := new(models.RepositoryPool)
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Select("size").Where("identity = ?", req.RepositoryIdentity).First(rp).Error; err != nil {
		return
	}
	ub := new(models.UserBasic)
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Select("now_volume", "total_volume").Where("identity = ?", userIdentity).First(ub).Error; err != nil {
		return
	}
	if ub.NowVolume+rp.Size > ub.TotalVolume {
		err = errors.New("storage capacity exceeded")
		return
	}

	// Update current capacity
	log.Printf("[UserRepositorySave] Updating user capacity: user=%s, file size=%d, current used=%d, after update=%d",
		userIdentity, rp.Size, ub.NowVolume, ub.NowVolume+rp.Size)
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		UpdateColumn("now_volume", gorm.Expr("now_volume + ?", rp.Size)).Error; err != nil {
		log.Printf("[UserRepositorySave] Failed to update capacity: %v", err)
		return
	}
	log.Printf("[UserRepositorySave] Capacity updated successfully")
	// Create association record
	ur := &models.UserRepository{
		Identity:           helper.UUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                req.Ext,
		Name:               req.Name,
	}
	err = l.svcCtx.DB.WithContext(l.ctx).Create(ur).Error
	return
}
