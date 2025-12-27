package logic

import (
	"context"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
)

type FriendRequestListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendRequestListLogic {
	return &FriendRequestListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendRequestListLogic) FriendRequestList(req *types.FriendRequestListRequest, userIdentity string) (resp *types.FriendRequestListReply, err error) {
	resp = new(types.FriendRequestListReply)
	resp.List = make([]*types.FriendRequestItem, 0)

	query := l.svcCtx.DB.WithContext(l.ctx).Table("friend_request").
		Select("friend_request.identity, friend_request.from_user_identity, friend_request.to_user_identity, " +
			"friend_request.status, friend_request.message, friend_request.created_at, " +
			"from_user.name as from_user_name, to_user.name as to_user_name").
		Joins("LEFT JOIN user_basic as from_user ON friend_request.from_user_identity = from_user.identity").
		Joins("LEFT JOIN user_basic as to_user ON friend_request.to_user_identity = to_user.identity")

	// Filter by type
	if req.Type == "sent" {
		query = query.Where("friend_request.from_user_identity = ?", userIdentity)
	} else if req.Type == "received" {
		query = query.Where("friend_request.to_user_identity = ?", userIdentity)
	} else {
		// all
		query = query.Where("friend_request.from_user_identity = ? OR friend_request.to_user_identity = ?", userIdentity, userIdentity)
	}

	query = query.Where("friend_request.deleted_at IS NULL").
		Order("friend_request.created_at DESC")

	var results []struct {
		Identity         string
		FromUserIdentity string
		ToUserIdentity   string
		FromUserName     string
		ToUserName       string
		Status           string
		Message          string
		CreatedAt        string
	}

	err = query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		// Format created_at
		createdAt := r.CreatedAt
		if createdAt == "" {
			createdAt = time.Now().Format(define.Datetime)
		}

		resp.List = append(resp.List, &types.FriendRequestItem{
			Identity:         r.Identity,
			FromUserIdentity: r.FromUserIdentity,
			ToUserIdentity:   r.ToUserIdentity,
			FromUserName:     r.FromUserName,
			ToUserName:       r.ToUserName,
			Status:           r.Status,
			Message:          r.Message,
			CreatedAt:        createdAt,
		})
	}

	return
}
