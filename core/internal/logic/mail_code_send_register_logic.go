package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	// 检查是否在最小发送间隔内（60秒），防止频繁发送
	codeTTL, _ := l.svcCtx.RDB.TTL(l.ctx, req.Email).Result()
	if codeTTL.Seconds() > 60 {
		return nil, errors.New("请等待60秒后再重新发送验证码")
	}
	// 获取验证码
	code := helper.RandCode()
	log.Printf("[MailCodeSend] 准备发送验证码到邮箱: %s, 验证码: %s", req.Email, code)

	// 存储验证码（覆盖旧的）
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))

	// 发送验证码
	err = helper.MailSendCode(req.Email, code)
	if err != nil {
		// 如果发送失败，删除已存储的验证码
		l.svcCtx.RDB.Del(l.ctx, req.Email)
		log.Printf("[MailCodeSend] 发送失败: %v", err)
		return nil, fmt.Errorf("发送验证码失败: %v", err)
	}
	log.Printf("[MailCodeSend] 验证码发送成功到: %s", req.Email)
	return
}
