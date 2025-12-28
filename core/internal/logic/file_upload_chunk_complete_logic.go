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

type FileUploadChunkCompleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkCompleteLogic {
	return &FileUploadChunkCompleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkCompleteLogic) FileUploadChunkComplete(req *types.FileUploadChunkCompleteRequest) (resp *types.FileUploadChunkCompleteReply, err error) {
	parts := make([]helper.MultipartPart, 0, len(req.Parts))
	for _, v := range req.Parts {
		parts = append(parts, helper.MultipartPart{
			ETag:       v.Etag,
			PartNumber: int32(v.PartNumber),
		})
	}
	if err = helper.S3PartUploadComplete(req.Key, req.UploadId, parts); err != nil {
		return
	}

	// Check if file already exists (deduplication)
	log.Printf("[FileUploadChunkComplete] Checking for existing file with xxHash64: %s (length: %d)", req.Md5, len(req.Md5))
	existingRp := new(models.RepositoryPool)
	err = l.svcCtx.DB.WithContext(l.ctx).Where("hash = ?", req.Md5).First(existingRp).Error
	if err == nil {
		log.Printf("[FileUploadChunkComplete] File already exists (deduplication): identity=%s, hash=%s", existingRp.Identity, existingRp.Hash)
		// File already exists, return existing identity
		resp = &types.FileUploadChunkCompleteReply{
			Identity: existingRp.Identity,
		}
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FileUploadChunkComplete] Database error while checking for existing file: %v", err)
		// Database error, not just "not found"
		return
	}
	log.Printf("[FileUploadChunkComplete] File not found, creating new record with xxHash64: %s", req.Md5)

	rp := &models.RepositoryPool{
		Identity: helper.UUID(),
		Hash:     req.Md5,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Key, // Store S3 key for permanent download endpoint
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(rp).Error; err != nil {
		return
	}

	resp = &types.FileUploadChunkCompleteReply{
		Identity: rp.Identity,
	}
	return
}
