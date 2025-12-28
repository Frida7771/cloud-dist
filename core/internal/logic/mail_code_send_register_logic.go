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

type MailCodeSendRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendRegisterLogic {
	return &MailCodeSendRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendRegisterLogic) MailCodeSendRegister(req *types.MailCodeSendRequest) (resp *types.MailCodeSendReply, err error) {
	// Check if email is not registered
	var cnt int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("email = ?", req.Email).Count(&cnt).Error; err != nil {
		return
	}
	if cnt > 0 {
		err = errors.New("email is already registered")
		return
	}
	// Check if within minimum send interval (60 seconds) to prevent frequent sending
	codeTTL, _ := l.svcCtx.RDB.TTL(l.ctx, req.Email).Result()
	if codeTTL.Seconds() > 60 {
		return nil, errors.New("please wait 60 seconds before resending verification code")
	}
	// Generate verification code
	code := helper.RandCode()
	log.Printf("[MailCodeSend] Preparing to send verification code to email: %s, code: %s", req.Email, code)

	// Store verification code (overwrite old one)
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))

	// Send verification code
	err = helper.MailSendCode(req.Email, code)
	if err != nil {
		// If sending fails, delete stored verification code
		l.svcCtx.RDB.Del(l.ctx, req.Email)
		log.Printf("[MailCodeSend] Send failed: %v", err)
		return nil, fmt.Errorf("failed to send verification code: %v", err)
	}
	log.Printf("[MailCodeSend] Verification code sent successfully to: %s", req.Email)
	return
}
