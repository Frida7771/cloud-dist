package logic

import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type FriendShareCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendShareCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendShareCreateLogic {
	return &FriendShareCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendShareCreateLogic) FriendShareCreate(req *types.FriendShareCreateRequest, fromUserIdentity string) (resp *types.FriendShareCreateReply, err error) {
	// Verify that users are friends
	var friendCount int64
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.Friend{}).
		Where("user_identity = ? AND friend_identity = ? AND status = ?", fromUserIdentity, req.ToUserIdentity, "active").
		Count(&friendCount).Error
	if err != nil {
		return nil, err
	}
	if friendCount == 0 {
		return nil, errors.New("users are not friends")
	}

	// Get user repository info
	var ur models.UserRepository
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ? AND user_identity = ?", req.UserRepositoryIdentity, fromUserIdentity).
		First(&ur).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("file not found")
	}
	if err != nil {
		return nil, err
	}

	// Create friend share
	fs := &models.FriendShare{
		Identity:               helper.UUID(),
		FromUserIdentity:       fromUserIdentity,
		ToUserIdentity:         req.ToUserIdentity,
		RepositoryIdentity:     ur.RepositoryIdentity,
		UserRepositoryIdentity: req.UserRepositoryIdentity,
		Message:                req.Message,
		IsRead:                 false,
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Create(fs).Error
	if err != nil {
		return nil, err
	}

	resp = &types.FriendShareCreateReply{
		Identity: fs.Identity,
	}
	return
}
