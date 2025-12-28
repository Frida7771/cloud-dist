package logic

import (
	"context"
	"errors"
	"log"

	"cloud-dist/core/helper"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
	"cloud-dist/core/svc"

	"gorm.io/gorm"
)

type ShareBasicSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicSaveLogic {
	return &ShareBasicSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicSaveLogic) ShareBasicSave(req *types.ShareBasicSaveRequest, userIdentity string) (resp *types.ShareBasicSaveReply, err error) {
	// Verify user identity is provided (should be checked in handler, but double-check here)
	if userIdentity == "" {
		return nil, errors.New("unauthorized: user identity is required")
	}

	// Verify the repository exists in repository_pool
	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).Where("identity = ?", req.RepositoryIdentity).First(rp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("resource does not exist")
	}
	if err != nil {
		return nil, err
	}

	// Check if file already exists in user's repository (global deduplication - user level)
	var existingUR models.UserRepository
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_identity = ?", userIdentity).
		Where("repository_identity = ?", req.RepositoryIdentity).
		Where("deleted_at IS NULL").
		First(&existingUR).Error
	if err == nil {
		// File already exists in user's repository (anywhere)
		return nil, errors.New("file already exists in your repository")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check user storage capacity
	ub := new(models.UserBasic)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Select("now_volume", "total_volume").
		Where("identity = ?", userIdentity).First(ub).Error
	if err != nil {
		return nil, err
	}

	if rp.Size+ub.NowVolume > ub.TotalVolume {
		return nil, errors.New("storage capacity exceeded")
	}

	// Create user repository entry
	ur := &models.UserRepository{
		Identity:           helper.UUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                rp.Ext,
		Name:               rp.Name,
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(ur).Error; err != nil {
		return nil, err
	}

	// Update user storage capacity
	log.Printf("[ShareBasicSave] Updating user capacity: user=%s, file size=%d, current used=%d, after update=%d",
		userIdentity, rp.Size, ub.NowVolume, ub.NowVolume+rp.Size)
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		UpdateColumn("now_volume", gorm.Expr("now_volume + ?", rp.Size)).Error; err != nil {
		log.Printf("[ShareBasicSave] Failed to update capacity: %v", err)
		return nil, err
	}
	log.Printf("[ShareBasicSave] Capacity updated successfully")

	resp = &types.ShareBasicSaveReply{Identity: ur.Identity}
	return
}
