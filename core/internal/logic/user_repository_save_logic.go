package logic

import (
	"context"
	"errors"
	"log"

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
	log.Printf("[UserRepositorySave] 更新用户容量: 用户=%s, 文件大小=%d, 当前已使用=%d, 更新后=%d",
		userIdentity, rp.Size, ub.NowVolume, ub.NowVolume+rp.Size)
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		UpdateColumn("now_volume", gorm.Expr("now_volume + ?", rp.Size)).Error; err != nil {
		log.Printf("[UserRepositorySave] 更新容量失败: %v", err)
		return
	}
	log.Printf("[UserRepositorySave] 容量更新成功")
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
