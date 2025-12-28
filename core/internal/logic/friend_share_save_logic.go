package logic

import (
	"context"
	"errors"

	"cloud-dist/core/helper"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
	"cloud-dist/core/svc"

	"gorm.io/gorm"
)

type FriendShareSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendShareSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendShareSaveLogic {
	return &FriendShareSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendShareSaveLogic) FriendShareSave(req *types.FriendShareSaveRequest, userIdentity string) (resp *types.FriendShareSaveReply, err error) {
	// First, verify the friend share record exists
	fs := new(models.FriendShare)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", req.ShareIdentity).
		First(fs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("share record not found")
		}
		return nil, err
	}

	// Verify that the user is the receiver (only receiver can save the file)
	if fs.ToUserIdentity != userIdentity {
		return nil, errors.New("access denied: only the receiver can save this file")
	}

	// Verify that both users are friends
	var friendCount int64
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.Friend{}).
		Where("((user_identity = ? AND friend_identity = ?) OR (user_identity = ? AND friend_identity = ?)) AND status = ?",
			fs.FromUserIdentity, fs.ToUserIdentity,
			fs.ToUserIdentity, fs.FromUserIdentity,
			"active").
		Count(&friendCount).Error
	if err != nil {
		return nil, err
	}

	if friendCount == 0 {
		return nil, errors.New("access denied: friendship relationship not found")
	}

	// Get file info from repository_pool
	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).Where("identity = ?", fs.RepositoryIdentity).First(rp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("resource does not exist")
	}
	if err != nil {
		return nil, err
	}

	// Check if file already exists in user's repository (global deduplication - user level, not folder level)
	var existingUR models.UserRepository
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_identity = ?", userIdentity).
		Where("repository_identity = ?", fs.RepositoryIdentity).
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
		RepositoryIdentity: fs.RepositoryIdentity,
		Ext:                rp.Ext,
		Name:               rp.Name,
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(ur).Error; err != nil {
		return nil, err
	}

	// Update user storage capacity
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		UpdateColumn("now_volume", gorm.Expr("now_volume + ?", rp.Size)).Error
	if err != nil {
		return nil, err
	}

	resp = &types.FriendShareSaveReply{Identity: ur.Identity}
	return
}
