package logic

import (
	"context"
	"time"

	"cloud-dist/core/define"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
)

type UserLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogoutLogic {
	return &UserLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogoutLogic) UserLogout(req *types.UserLogoutRequest, token string) (resp *types.UserLogoutReply, err error) {
	// Add token to blacklist in Redis
	// Store with expiration time equal to token's remaining validity
	// Use token as key, value can be "blacklisted" or timestamp
	// Expiration should match token's remaining TTL (max TokenExpire seconds)
	
	// Calculate expiration: use TokenExpire as max expiration
	// In practice, you might want to parse the token to get its actual expiration
	// For simplicity, we'll use TokenExpire as the blacklist TTL
	blacklistKey := "token:blacklist:" + token
	expiration := time.Duration(define.TokenExpire) * time.Second
	
	err = l.svcCtx.RDB.Set(l.ctx, blacklistKey, "blacklisted", expiration).Err()
	if err != nil {
		return nil, err
	}
	
	resp = &types.UserLogoutReply{}
	return
}


