package logic

import (
	"context"
	"errors"

	"cloud-disk/core/svc"
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
	var parentID int64 = 0 // Default to root (parent_id = 0)
	
	// If parent_identity is provided (not empty), find the parent folder
	if req.ParentIdnetity != "" {
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
		parentID = parentData.ID
	}
	// If parent_identity is empty, parentID remains 0 (root directory)

	err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Where("identity = ? AND user_identity = ?", req.Idnetity, userIdentity).
		Update("parent_id", parentID).Error
	return
}
