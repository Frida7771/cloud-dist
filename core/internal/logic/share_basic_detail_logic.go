package logic

import (
	"context"

	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type ShareBasicDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDetailLogic {
	return &ShareBasicDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDetailLogic) ShareBasicDetail(req *types.ShareBasicDetailRequest) (resp *types.ShareBasicDetailReply, err error) {
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.ShareBasic{}).
		Where("identity = ?", req.Identity).
		UpdateColumn("click_num", gorm.Expr("click_num + 1")).Error; err != nil {
		return
	}

	resp = new(types.ShareBasicDetailReply)
	err = l.svcCtx.DB.WithContext(l.ctx).Table("share_basic").
		Select("share_basic.repository_identity, user_repository.name, repository_pool.ext, repository_pool.size, repository_pool.path").
		Joins("LEFT JOIN repository_pool ON share_basic.repository_identity = repository_pool.identity").
		Joins("LEFT JOIN user_repository ON user_repository.identity = share_basic.user_repository_identity").
		Where("share_basic.identity = ?", req.Identity).
		Take(resp).Error
	if err != nil {
		return
	}

	// Convert repository identity to download endpoint URL
	// This provides permanent download links that don't expire
	if resp.RepositoryIdentity != "" {
		resp.Path = "/file/download?identity=" + resp.RepositoryIdentity
	}

	return
}
