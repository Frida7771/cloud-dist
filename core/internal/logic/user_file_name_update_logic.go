package logic

import (
	"context"
	"errors"

	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
)

type UserFileNameUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameUpdateLogic {
	return &UserFileNameUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameUpdateLogic) UserFileNameUpdate(req *types.UserFileNameUpdateRequest, userIdentity string) (resp *types.UserFileNameUpdateReply, err error) {
	parentQuery := l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Select("parent_id").Where("identity = ?", req.Identity).Limit(1)

	var cnt int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Where("name = ?", req.Name).
		Where("parent_id = (?)", parentQuery).
		Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("name already exists")
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).
		Update("name", req.Name).Error
	return
}
