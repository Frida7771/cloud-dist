package logic

import (
	"context"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type UserFileDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileDeleteLogic {
	return &UserFileDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteRequest, userIdentity string) (resp *types.UserFileDeleteReply, err error) {
	var size int64
	err = l.svcCtx.DB.WithContext(l.ctx).Table("repository_pool").
		Select("repository_pool.size").
		Joins("JOIN user_repository ON user_repository.repository_identity = repository_pool.identity").
		Where("user_repository.identity = ?", req.Identity).
		Limit(1).
		Scan(&size).Error
	if err != nil {
		return
	}

	if size > 0 {
		if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
			Where("identity = ?", userIdentity).
			UpdateColumn("now_volume", gorm.Expr("now_volume - ?", size)).Error; err != nil {
			return
		}
	}

	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).
		Delete(&models.UserRepository{}).Error
	return
}
