package logic

import (
	"context"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
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

func (l *RefreshAuthorizationLogic) RefreshAuthorization(req *types.RefreshAuthorizationRequest, authorization string) (resp *types.RefreshAuthorizationReply, err error) {
	uc, err := helper.AnalyzeToken(authorization)
	if err != nil {
		return
	}
	token, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.TokenExpire)
	if err != nil {
		return
	}
	refreshToken, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.RefreshTokenExpire)
	if err != nil {
		return
	}

	resp = new(types.RefreshAuthorizationReply)
	resp.Token = token
	resp.RefreshToken = refreshToken
	return
}
