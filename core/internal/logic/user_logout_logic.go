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

func (l *UserLogoutLogic) UserLogout(req *types.UserLogoutRequest, refreshToken string) (resp *types.UserLogoutReply, err error) {
	// Add refresh token to blacklist in Redis
	// Store with expiration time equal to refresh token's remaining validity
	// Use refresh token as key, value can be "blacklisted" or timestamp
	// Expiration should match refresh token's remaining TTL (max RefreshTokenExpire seconds)
	
	// Calculate expiration: use RefreshTokenExpire as max expiration
	// In practice, you might want to parse the token to get its actual expiration
	// For simplicity, we'll use RefreshTokenExpire as the blacklist TTL
	blacklistKey := "refresh_token:blacklist:" + refreshToken
	expiration := time.Duration(define.RefreshTokenExpire) * time.Second
	
	err = l.svcCtx.RDB.Set(l.ctx, blacklistKey, "blacklisted", expiration).Err()
	if err != nil {
		return nil, err
	}
	
	resp = &types.UserLogoutReply{}
	return
}


