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

type UserRepositorySaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRepositorySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRepositorySaveLogic {
	return &UserRepositorySaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRepositorySaveLogic) UserRepositorySave(req *types.UserRepositorySaveRequest, userIdentity string) (resp *types.UserRepositorySaveReply, err error) {
	// 判断文件是否超容量
	rp := new(models.RepositoryPool)
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Select("size").Where("identity = ?", req.RepositoryIdentity).First(rp).Error; err != nil {
		return
	}
	ub := new(models.UserBasic)
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Select("now_volume", "total_volume").Where("identity = ?", userIdentity).First(ub).Error; err != nil {
		return
	}
	if ub.NowVolume+rp.Size > ub.TotalVolume {
		err = errors.New("已超出当前容量")
		return
	}

	// 更新当前容量
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		UpdateColumn("now_volume", gorm.Expr("now_volume + ?", rp.Size)).Error; err != nil {
		return
	}
	// 新增关联记录
	ur := &models.UserRepository{
		Identity:           helper.UUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                req.Ext,
		Name:               req.Name,
	}
	err = l.svcCtx.DB.WithContext(l.ctx).Create(ur).Error
	return
}
