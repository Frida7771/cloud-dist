package logic

import (
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type UserFileMoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveRequest, userIdentity string) (resp *types.UserFileMoveReply, err error) {
	//parentID
	parentData := new(models.UserRepository)
	err = l.svcCtx.DB.WithContext(l.ctx).
		Where("identity = ? AND user_identity = ?", req.ParentIdnetity, userIdentity).
		First(parentData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("folder does not exist")
	}
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Where("identity = ?", req.Idnetity).
		Update("parent_id", parentData.ID).Error
	return
}
