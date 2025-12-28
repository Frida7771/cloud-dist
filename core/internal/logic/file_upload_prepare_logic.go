package logic

import (
	"context"
	"errors"
	"log"

	"cloud-dist/core/helper"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
	"cloud-dist/core/svc"

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

	// First check if file already exists in repository (deduplication)
	rp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).Where("hash = ?", req.Md5).First(rp).Error
	if err == nil {
		log.Printf("[FileUploadPrepare] File already exists (deduplication): identity=%s, hash=%s", rp.Identity, rp.Hash)
		resp.Identity = rp.Identity
		return resp, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FileUploadPrepare] Database error while checking for existing file: %v", err)
		return nil, err
	}

	// Create new upload task
	// Note: Resume functionality is handled by frontend localStorage
	log.Printf("[FileUploadPrepare] Creating new upload task")
	key, uploadId, err := helper.S3InitPart(req.Ext)
	if err != nil {
		return nil, err
	}
	resp.Key = key
	resp.UploadId = uploadId

	return
}
