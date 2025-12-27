package logic

import (
	"context"
	"errors"

	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type UserFolderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderListLogic {
	return &UserFolderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderListLogic) UserFolderList(req *types.UserFolderListRequest) (resp *types.UserFolderListReply, err error) {
	resp = new(types.UserFolderListReply)
	folders := make([]*types.UserFolder, 0)

	var parentID int64
	if req.Identity != "" {
		ur := new(models.UserRepository)
		err = l.svcCtx.DB.WithContext(l.ctx).Select("id").
			Where("identity = ?", req.Identity).First(ur).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if err == nil {
			parentID = ur.ID
		}
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Table("user_repository").
		Select("identity, name").
		Where("parent_id = ?", parentID).
		Where("deleted_at IS NULL").
		Find(&folders).Error
	if err != nil {
		return
	}

	resp.List = folders
	return
}
