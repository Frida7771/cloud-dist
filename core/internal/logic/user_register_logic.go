package logic

import (
	"context"
	"errors"
	"log"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"cloud-disk/core/svc"
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
	// Verify code
	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, errors.New("verification code is empty for this email")
	}
	if code != req.Code {
		err = errors.New("verification code is incorrect")
		return
	}
	var cnt int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("name = ?", req.Name).Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt > 0 {
		err = errors.New("username already exists")
		return
	}
	// Save user data
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user := &models.UserBasic{
		Identity:    helper.UUID(),
		Name:        req.Name,
		Password:    hashedPassword,
		Email:       req.Email,
		NowVolume:   0,           // Initial used volume is 0
		TotalVolume: 16106127360, // Default total volume 15GB (15 * 1024 * 1024 * 1024)
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(user).Error; err != nil {
		return nil, err
	}
	log.Println("insert user row:", user.ID)
	return
}
