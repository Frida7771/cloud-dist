package logic

import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type ShareBasicCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicCreateLogic {
	return &ShareBasicCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicCreateLogic) ShareBasicCreate(req *types.ShareBasicCreateRequest, userIdentity string) (resp *types.ShareBasicCreateReply, err error) {
	uuid := helper.UUID()
	ur := new(models.UserRepository)
	err = l.svcCtx.DB.WithContext(l.ctx).Where("identity = ?", req.UserRepositoryIdentity).First(ur).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user repository not found")
	}
	if err != nil {
		return nil, err
	}

	data := &models.ShareBasic{
		Identity:               uuid,
		UserIdentity:           userIdentity,
		UserRepositoryIdentity: req.UserRepositoryIdentity,
		RepositoryIdentity:     ur.RepositoryIdentity,
		ExpiredTime:            req.ExpiredTime,
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(data).Error; err != nil {
		return
	}
	resp = &types.ShareBasicCreateReply{
		Identity: uuid,
	}
	return
}
