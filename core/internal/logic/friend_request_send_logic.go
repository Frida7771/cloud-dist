package logic

import (
	"context"
	"errors"
	"strings"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type FriendRequestSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendRequestSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendRequestSendLogic {
	return &FriendRequestSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendRequestSendLogic) FriendRequestSend(req *types.FriendRequestSendRequest, fromUserIdentity string) (resp *types.FriendRequestSendReply, err error) {
	// Find target user by email or identity
	var toUser models.UserBasic
	toUserIdentity := strings.TrimSpace(req.ToUserIdentity)

	// Check if it's an email or user identity
	if strings.Contains(toUserIdentity, "@") {
		// It's an email
		err = l.svcCtx.DB.WithContext(l.ctx).
			Where("email = ?", toUserIdentity).
			First(&toUser).Error
	} else {
		// It's a user identity
		err = l.svcCtx.DB.WithContext(l.ctx).
			Where("identity = ?", toUserIdentity).
			First(&toUser).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Can't send request to yourself
	if toUser.Identity == fromUserIdentity {
		return nil, errors.New("cannot send friend request to yourself")
	}

	// Check if already friends
	var friendCount int64
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.Friend{}).
		Where("user_identity = ? AND friend_identity = ?", fromUserIdentity, toUser.Identity).
		Count(&friendCount).Error
	if err != nil {
		return nil, err
	}
	if friendCount > 0 {
		return nil, errors.New("already friends")
	}

	// Check if there's a pending request
	var existingRequest models.FriendRequest
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("((from_user_identity = ? AND to_user_identity = ?) OR (from_user_identity = ? AND to_user_identity = ?)) AND status = ?",
			fromUserIdentity, toUser.Identity, toUser.Identity, fromUserIdentity, "pending").
		First(&existingRequest).Error
	if err == nil {
		return nil, errors.New("friend request already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create friend request
	fr := &models.FriendRequest{
		Identity:         helper.UUID(),
		FromUserIdentity: fromUserIdentity,
		ToUserIdentity:   toUser.Identity,
		Status:           "pending",
		Message:          req.Message,
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Create(fr).Error
	if err != nil {
		return nil, err
	}

	resp = &types.FriendRequestSendReply{
		Identity: fr.Identity,
	}
	return
}
