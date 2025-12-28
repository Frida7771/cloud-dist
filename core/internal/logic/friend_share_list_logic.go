package logic

import (
	"context"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/svc"
)

type FriendShareListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendShareListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendShareListLogic {
	return &FriendShareListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendShareListLogic) FriendShareList(req *types.FriendShareListRequest, userIdentity string) (resp *types.FriendShareListReply, err error) {
	resp = new(types.FriendShareListReply)
	resp.List = make([]*types.FriendShareItem, 0)

	query := l.svcCtx.DB.WithContext(l.ctx).Table("friend_share").
		Select("friend_share.identity, friend_share.from_user_identity, friend_share.to_user_identity, " +
			"friend_share.repository_identity, friend_share.user_repository_identity, " +
			"friend_share.message, friend_share.is_read, friend_share.created_at, " +
			"from_user.name as from_user_name, to_user.name as to_user_name, " +
			"repository_pool.name as file_name, repository_pool.ext as file_ext, repository_pool.size as file_size, " +
			"repository_pool.path as s3_key").
		Joins("LEFT JOIN user_basic as from_user ON friend_share.from_user_identity = from_user.identity").
		Joins("LEFT JOIN user_basic as to_user ON friend_share.to_user_identity = to_user.identity").
		Joins("LEFT JOIN repository_pool ON friend_share.repository_identity = repository_pool.identity")

	// Filter by type
	if req.Type == "sent" {
		query = query.Where("friend_share.from_user_identity = ?", userIdentity)
	} else if req.Type == "received" {
		query = query.Where("friend_share.to_user_identity = ?", userIdentity)
	} else {
		// all
		query = query.Where("friend_share.from_user_identity = ? OR friend_share.to_user_identity = ?", userIdentity, userIdentity)
	}

	query = query.Where("friend_share.deleted_at IS NULL").
		Order("friend_share.created_at DESC")

	var results []struct {
		Identity               string
		FromUserIdentity       string
		FromUserName           string
		ToUserIdentity         string
		ToUserName             string
		RepositoryIdentity     string
		UserRepositoryIdentity string
		FileName               string
		FileExt                string
		FileSize               int64
		Message                string
		IsRead                 bool
		CreatedAt              string
		S3Key                  string
	}

	err = query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		item := &types.FriendShareItem{
			Identity:               r.Identity,
			FromUserIdentity:       r.FromUserIdentity,
			FromUserName:           r.FromUserName,
			ToUserIdentity:         r.ToUserIdentity,
			ToUserName:             r.ToUserName,
			RepositoryIdentity:     r.RepositoryIdentity,
			UserRepositoryIdentity: r.UserRepositoryIdentity,
			FileName:               r.FileName,
			FileExt:                r.FileExt,
			FileSize:               r.FileSize,
			Message:                r.Message,
			IsRead:                 r.IsRead,
			CreatedAt:              r.CreatedAt,
		}

		// Use friend share download endpoint (no expiration, verifies friendship)
		// This endpoint verifies that both users are friends and allows permanent access
		if r.Identity != "" {
			item.Path = "/friend/share/download?identity=" + r.Identity
		} else if r.RepositoryIdentity != "" {
			// Fallback to repository identity if share identity is not available
			item.Path = "/friend/share/download?identity=" + r.RepositoryIdentity
		}

		// Format created_at
		if item.CreatedAt == "" {
			item.CreatedAt = time.Now().Format(define.Datetime)
		}

		resp.List = append(resp.List, item)
	}

	return
}
