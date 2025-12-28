package logic

import (
	"context"
	"errors"

	"cloud-dist/core/helper"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"

	"gorm.io/gorm"
)

type UserPasswordResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPasswordResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPasswordResetLogic {
	return &UserPasswordResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPasswordResetLogic) UserPasswordReset(req *types.UserPasswordResetRequest) (resp *types.UserPasswordResetReply, err error) {
	// Get user by email
	user := new(models.UserBasic)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("email = ?", req.Email).
		First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("email not found")
	}
	if err != nil {
		return nil, err
	}

	// Verify email verification code
	redisKey := "password_reset:" + req.Email
	storedCode, err := l.svcCtx.RDB.Get(l.ctx, redisKey).Result()
	if err != nil {
		return nil, errors.New("verification code is expired or not found, please request a new code")
	}
	if storedCode != req.Code {
		return nil, errors.New("verification code is incorrect")
	}

	// Validate new password
	if req.NewPassword == "" {
		return nil, errors.New("new password cannot be empty")
	}
	if len(req.NewPassword) < 6 {
		return nil, errors.New("new password must be at least 6 characters")
	}

	// Hash new password using bcrypt
	newPasswordHash, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errors.New("failed to hash new password")
	}

	// Update password
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("email = ?", req.Email).
		Update("password", newPasswordHash).Error
	if err != nil {
		return nil, err
	}

	// Delete verification code after successful password reset
	l.svcCtx.RDB.Del(l.ctx, redisKey)

	resp = &types.UserPasswordResetReply{}
	return
}

