package logic

import (
	"context"
	"errors"
	"log"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterReply, err error) {
	// 判断code是否一致
	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, errors.New("该邮箱的验证码为空")
	}
	if code != req.Code {
		err = errors.New("验证码错误")
		return
	}
	var cnt int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("name = ?", req.Name).Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt > 0 {
		err = errors.New("用户名已存在")
		return
	}
	// 数据入库
	user := &models.UserBasic{
		Identity:    helper.UUID(),
		Name:        req.Name,
		Password:    helper.Md5(req.Password),
		Email:       req.Email,
		NowVolume:   0,          // 初始已使用容量为 0
		TotalVolume: 5368709120, // 默认总容量 5GB (5 * 1024 * 1024 * 1024)
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(user).Error; err != nil {
		return nil, err
	}
	log.Println("insert user row:", user.ID)
	return
}
