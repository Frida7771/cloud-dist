package logic

import (
	"context"
	"errors"

	"cloud-dist/core/helper"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"

	"gorm.io/gorm"
)

type FriendRequestRespondLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendRequestRespondLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendRequestRespondLogic {
	return &FriendRequestRespondLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendRequestRespondLogic) FriendRequestRespond(req *types.FriendRequestRespondRequest, userIdentity string) (resp *types.FriendRequestRespondReply, err error) {
	// Get the friend request
	var fr models.FriendRequest
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ? AND to_user_identity = ?", req.Identity, userIdentity).
		First(&fr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("friend request not found")
	}
	if err != nil {
		return nil, err
	}

	// Check if already processed
	if fr.Status != "pending" {
		return nil, errors.New("friend request already processed")
	}

	// Update request status
	newStatus := req.Action
	if newStatus != "accept" && newStatus != "reject" {
		return nil, errors.New("invalid action, must be 'accept' or 'reject'")
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.FriendRequest{}).
		Where("identity = ?", req.Identity).
		Update("status", newStatus).Error
	if err != nil {
		return nil, err
	}

	// If accepted, create friendship (bidirectional)
	if newStatus == "accept" {
		// Create friend relationship: user -> friend
		friend1 := &models.Friend{
			Identity:       helper.UUID(),
			UserIdentity:   fr.ToUserIdentity,   // The user who accepted
			FriendIdentity: fr.FromUserIdentity, // The user who sent the request
			Status:         "active",
		}
		err = l.svcCtx.DB.WithContext(l.ctx).Create(friend1).Error
		if err != nil {
			return nil, err
		}

		// Create reverse friend relationship: friend -> user
		friend2 := &models.Friend{
			Identity:       helper.UUID(),
			UserIdentity:   fr.FromUserIdentity, // The user who sent the request
			FriendIdentity: fr.ToUserIdentity,   // The user who accepted
			Status:         "active",
		}
		err = l.svcCtx.DB.WithContext(l.ctx).Create(friend2).Error
		if err != nil {
			return nil, err
		}
	}

	resp = &types.FriendRequestRespondReply{}
	return
}
