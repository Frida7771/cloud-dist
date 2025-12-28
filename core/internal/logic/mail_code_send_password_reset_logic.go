package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"cloud-disk/core/svc"

	"gorm.io/gorm"
)

type MailCodeSendPasswordResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendPasswordResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendPasswordResetLogic {
	return &MailCodeSendPasswordResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendPasswordResetLogic) MailCodeSendPasswordReset(req *types.MailCodeSendPasswordResetRequest) (resp *types.MailCodeSendPasswordResetReply, err error) {
	// Check if email is registered
	var user models.UserBasic
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("email = ?", req.Email).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("email not found")
	}
	if err != nil {
		return nil, err
	}

	// Check if within minimum send interval (60 seconds) to prevent frequent sending
	redisKey := "password_reset:" + req.Email
	codeTTL, _ := l.svcCtx.RDB.TTL(l.ctx, redisKey).Result()
	if codeTTL.Seconds() > 60 {
		return nil, errors.New("please wait 60 seconds before resending verification code")
	}

	// Generate verification code
	code := helper.RandCode()
	log.Printf("[MailCodeSendPasswordReset] Preparing to send verification code to email: %s, code: %s", req.Email, code)

	// Store verification code (overwrite old one)
	l.svcCtx.RDB.Set(l.ctx, redisKey, code, time.Second*time.Duration(define.CodeExpire))

	// Send verification code
	err = helper.MailSendCode(req.Email, code)
	if err != nil {
		// If sending fails, delete stored verification code
		l.svcCtx.RDB.Del(l.ctx, redisKey)
		log.Printf("[MailCodeSendPasswordReset] Send failed: %v", err)
		return nil, fmt.Errorf("failed to send verification code: %v", err)
	}
	log.Printf("[MailCodeSendPasswordReset] Verification code sent successfully to: %s", req.Email)
	return
}
