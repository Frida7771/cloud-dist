package logic

import (
	"context"
	"errors"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type UserLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginReply, err error) {
	user := new(models.UserBasic)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("name = ?", req.Name).
		First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("username or password is incorrect")
	}
	if err != nil {
		return nil, err
	}
	
	// Verify password using bcrypt
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("username or password is incorrect")
	}

	token, err := helper.GenerateToken(int(user.ID), user.Identity, user.Name, define.TokenExpire)
	if err != nil {
		return nil, err
	}
	refreshToken, err := helper.GenerateToken(int(user.ID), user.Identity, user.Name, define.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}

	resp = new(types.LoginReply)
	resp.Token = token
	resp.RefreshToken = refreshToken
	return
}
