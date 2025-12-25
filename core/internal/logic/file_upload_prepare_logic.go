package logic

import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type FileUploadPrepareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadPrepareLogic {
	return &FileUploadPrepareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadPrepareLogic) FileUploadPrepare(req *types.FileUploadPrepareRequest) (resp *types.FileUploadPrepareReply, err error) {
	resp = new(types.FileUploadPrepareReply)

	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).Where("hash = ?", req.Md5).First(rp).Error
	if err == nil {
		resp.Identity = rp.Identity
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	key, uploadId, err := helper.S3InitPart(req.Ext)
	if err != nil {
		return nil, err
	}
	resp.Key = key
	resp.UploadId = uploadId
	return
}
