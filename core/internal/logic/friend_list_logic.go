package logic

import (
	"context"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListRequest, userIdentity string) (resp *types.FriendListReply, err error) {
	resp = new(types.FriendListReply)
	resp.List = make([]*types.FriendItem, 0)

	var results []struct {
		Identity     string
		UserIdentity string
		UserName     string
		UserEmail    string
		Status       string
		CreatedAt    string
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Table("friend").
		Select("friend.identity, friend.friend_identity as user_identity, "+
			"user_basic.name as user_name, user_basic.email as user_email, "+
			"friend.status, friend.created_at").
		Joins("LEFT JOIN user_basic ON friend.friend_identity = user_basic.identity").
		Where("friend.user_identity = ?", userIdentity).
		Where("friend.status = ?", "active").
		Where("friend.deleted_at IS NULL").
		Order("friend.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, r := range results {
		// Format created_at
		createdAt := r.CreatedAt
		if createdAt == "" {
			createdAt = time.Now().Format(define.Datetime)
		}

		resp.List = append(resp.List, &types.FriendItem{
			Identity:     r.Identity,
			UserIdentity: r.UserIdentity,
			UserName:     r.UserName,
			UserEmail:    r.UserEmail,
			Status:       r.Status,
			CreatedAt:    createdAt,
		})
	}

	return
}
