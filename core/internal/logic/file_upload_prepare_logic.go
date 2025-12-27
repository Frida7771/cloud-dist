package logic

import (
	"context"
	"errors"
	"log"

	"cloud-disk/core/helper"
	"cloud-disk/core/svc"
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

	// Log the hash being searched (for debugging)
	log.Printf("[FileUploadPrepare] Checking for existing file with xxHash64: %s (length: %d)", req.Md5, len(req.Md5))

	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).Where("hash = ?", req.Md5).First(rp).Error
	if err == nil {
		log.Printf("[FileUploadPrepare] File already exists (deduplication): identity=%s, hash=%s", rp.Identity, rp.Hash)
		resp.Identity = rp.Identity
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FileUploadPrepare] Database error while checking for existing file: %v", err)
		return nil, err
	}
	log.Printf("[FileUploadPrepare] File not found, proceeding with new upload")

	key, uploadId, err := helper.S3InitPart(req.Ext)
	if err != nil {
		return nil, err
	}
	resp.Key = key
	resp.UploadId = uploadId
	return
}
