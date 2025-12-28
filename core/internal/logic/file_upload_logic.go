package logic

import (
	"context"
	"log"

	"cloud-dist/core/helper"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
)

type FileUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest) (resp *types.FileUploadReply, err error) {
	log.Printf("[FileUploadLogic] Creating file record: filename=%s, ext=%s, size=%d, xxHash64=%s", req.Name, req.Ext, req.Size, req.Hash)

	rp := &models.RepositoryPool{
		Identity: helper.UUID(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(rp).Error; err != nil {
		log.Printf("[FileUploadLogic] Failed to save file record: %v", err)
		return nil, err
	}
	log.Printf("[FileUploadLogic] File record saved successfully: identity=%s", rp.Identity)

	resp = &types.FileUploadReply{
		Identity: rp.Identity,
		Ext:      rp.Ext,
		Name:     rp.Name,
	}
	return
}
