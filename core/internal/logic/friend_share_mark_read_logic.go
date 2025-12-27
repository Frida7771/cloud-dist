package logic

import (
	"context"
	"errors"

	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type FriendShareMarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendShareMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendShareMarkReadLogic {
	return &FriendShareMarkReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendShareMarkReadLogic) FriendShareMarkRead(req *types.FriendShareMarkReadRequest, userIdentity string) (resp *types.FriendShareMarkReadReply, err error) {
	// Verify that the share belongs to the user
	var fs models.FriendShare
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ? AND to_user_identity = ?", req.Identity, userIdentity).
		First(&fs).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("share not found")
	}
	if err != nil {
		return nil, err
	}

	// Mark as read
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.FriendShare{}).
		Where("identity = ?", req.Identity).
		Update("is_read", true).Error
	if err != nil {
		return nil, err
	}

	resp = &types.FriendShareMarkReadReply{}
	return
}
