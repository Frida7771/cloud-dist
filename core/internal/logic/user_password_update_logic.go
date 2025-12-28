package logic

import (
	"context"
	"errors"

	"cloud-dist/core/helper"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
	"cloud-dist/core/svc"

	"gorm.io/gorm"
)

type UserPasswordUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPasswordUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPasswordUpdateLogic {
	return &UserPasswordUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPasswordUpdateLogic) UserPasswordUpdate(req *types.UserPasswordUpdateRequest, userIdentity string) (resp *types.UserPasswordUpdateReply, err error) {
	// Get user by identity
	user := new(models.UserBasic)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", userIdentity).
		First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Verify email verification code
	redisKey := "password_update:" + user.Email
	storedCode, err := l.svcCtx.RDB.Get(l.ctx, redisKey).Result()
	if err != nil {
		return nil, errors.New("verification code is expired or not found, please request a new code")
	}
	if storedCode != req.Code {
		return nil, errors.New("verification code is incorrect")
	}

	// Verify old password using bcrypt
	if !helper.CheckPasswordHash(req.OldPassword, user.Password) {
		return nil, errors.New("old password is incorrect")
	}

	// Validate new password
	if req.NewPassword == "" {
		return nil, errors.New("new password cannot be empty")
	}
	if len(req.NewPassword) < 6 {
		return nil, errors.New("new password must be at least 6 characters")
	}

	// Check if new password is same as old password
	if helper.CheckPasswordHash(req.NewPassword, user.Password) {
		return nil, errors.New("new password must be different from old password")
	}

	// Hash new password using bcrypt
	newPasswordHash, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errors.New("failed to hash new password")
	}

	// Update password
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		Update("password", newPasswordHash).Error
	if err != nil {
		return nil, err
	}

	// Delete verification code after successful password update
	l.svcCtx.RDB.Del(l.ctx, redisKey)

	resp = &types.UserPasswordUpdateReply{}
	return
}
