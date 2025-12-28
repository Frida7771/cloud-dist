package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud-dist/core/define"
	"cloud-dist/core/helper"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
)

type MailCodeSendPasswordUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendPasswordUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendPasswordUpdateLogic {
	return &MailCodeSendPasswordUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendPasswordUpdateLogic) MailCodeSendPasswordUpdate(req *types.MailCodeSendPasswordUpdateRequest, userIdentity string) (resp *types.MailCodeSendPasswordUpdateReply, err error) {
	// Get user by identity to verify email
	user := new(models.UserBasic)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ?", userIdentity).
		First(user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify that the email matches the user's email
	if user.Email != req.Email {
		return nil, errors.New("email does not match your account")
	}

	// Check if within minimum send interval (60 seconds) to prevent frequent sending
	redisKey := "password_update:" + req.Email
	codeTTL, _ := l.svcCtx.RDB.TTL(l.ctx, redisKey).Result()
	if codeTTL.Seconds() > 60 {
		return nil, errors.New("please wait 60 seconds before resending verification code")
	}

	// Generate verification code
	code := helper.RandCode()
	log.Printf("[MailCodeSendPasswordUpdate] Preparing to send verification code to email: %s, code: %s", req.Email, code)

	// Store verification code (overwrite old one)
	l.svcCtx.RDB.Set(l.ctx, redisKey, code, time.Second*time.Duration(define.CodeExpire))

	// Send verification code
	err = helper.MailSendCode(req.Email, code)
	if err != nil {
		// If sending fails, delete stored verification code
		l.svcCtx.RDB.Del(l.ctx, redisKey)
		log.Printf("[MailCodeSendPasswordUpdate] Send failed: %v", err)
		return nil, fmt.Errorf("failed to send verification code: %v", err)
	}
	log.Printf("[MailCodeSendPasswordUpdate] Verification code sent successfully to: %s", req.Email)
	return
}


