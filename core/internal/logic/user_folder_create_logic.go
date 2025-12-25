package logic

import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
)

type UserFolderCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderCreateLogic {
	return &UserFolderCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderCreateLogic) UserFolderCreate(req *types.UserFolderCreateRequest, userIdentity string) (resp *types.UserFolderCreateReply, err error) {
	// 判断当前名称在该层级下是否存在
	var cnt int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Where("name = ? AND parent_id = ? AND user_identity = ?", req.Name, req.ParentId, userIdentity).
		Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("该名称已存在")
	}
	// 创建文件夹
	data := &models.UserRepository{
		Identity:     helper.UUID(),
		UserIdentity: userIdentity,
		ParentId:     req.ParentId,
		Name:         req.Name,
	}
	err = l.svcCtx.DB.WithContext(l.ctx).Create(data).Error
	return
}
