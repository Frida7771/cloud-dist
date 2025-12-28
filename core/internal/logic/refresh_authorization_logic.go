package logic

import (
	"context"
	"errors"
	"time"

	"cloud-dist/core/define"
	"cloud-dist/core/helper"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
)

type RefreshAuthorizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshAuthorizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshAuthorizationLogic {
	return &RefreshAuthorizationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshAuthorizationLogic) RefreshAuthorization(req *types.RefreshAuthorizationRequest, refreshToken string) (resp *types.RefreshAuthorizationReply, err error) {
	// Parse refresh token
	uc, err := helper.AnalyzeToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if refresh token is blacklisted
	if l.svcCtx.RDB != nil {
		blacklistKey := "refresh_token:blacklist:" + refreshToken
		_, err := l.svcCtx.RDB.Get(l.ctx, blacklistKey).Result()
		if err == nil {
			// Refresh token found in blacklist
			return nil, errors.New("refresh token has been revoked")
		}
	}

	// Generate new access token
	token, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.TokenExpire)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token (rotate refresh token for security)
	newRefreshToken, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}

	// Optionally: blacklist the old refresh token (token rotation)
	// This ensures that if an old refresh token is stolen, it can't be used after refresh
	if l.svcCtx.RDB != nil {
		oldBlacklistKey := "refresh_token:blacklist:" + refreshToken
		expiration := time.Duration(define.RefreshTokenExpire) * time.Second
		l.svcCtx.RDB.Set(l.ctx, oldBlacklistKey, "blacklisted", expiration)
	}

	resp = new(types.RefreshAuthorizationReply)
	resp.Token = token
	resp.RefreshToken = newRefreshToken
	return
}
