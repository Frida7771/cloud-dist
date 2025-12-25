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
