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
	// Check if name already exists at this level
	var cnt int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Where("name = ? AND parent_id = ? AND user_identity = ?", req.Name, req.ParentId, userIdentity).
		Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("name already exists")
	}
	// Create folder
	data := &models.UserRepository{
		Identity:     helper.UUID(),
		UserIdentity: userIdentity,
		ParentId:     req.ParentId,
		Name:         req.Name,
		Ext:          "", // Empty ext indicates it's a folder
	}
	err = l.svcCtx.DB.WithContext(l.ctx).Create(data).Error
	if err == nil {
		resp = &types.UserFolderCreateReply{
			Identity: data.Identity,
		}
	}
	return
}
