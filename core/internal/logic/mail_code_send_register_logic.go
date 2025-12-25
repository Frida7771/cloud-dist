package logic

import (
	"context"
	"errors"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
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
	// 该邮箱未被注册
	var cnt int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserBasic{}).
		Where("email = ?", req.Email).Count(&cnt).Error; err != nil {
		return
	}
	if cnt > 0 {
		err = errors.New("该邮箱已被注册")
		return
	}
	// 校验验证码是否在有效期内
	codeTTL, _ := l.svcCtx.RDB.TTL(l.ctx, req.Email).Result()
	if codeTTL.Seconds() > 0 || codeTTL.Seconds() == -1 {
		return nil, errors.New("the verify code has not expired")
	}
	// 获取验证码
	code := helper.RandCode()
	// 存储验证码
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))
	// 发送验证码
	err = helper.MailSendCode(req.Email, code)
	return
}
