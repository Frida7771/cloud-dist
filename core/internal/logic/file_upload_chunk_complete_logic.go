package logic

import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
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

	rp := &models.RepositoryPool{
		Identity: helper.UUID(),
		Hash:     req.Md5,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     helper.S3ObjectURL(req.Key),
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(rp).Error; err != nil {
		return
	}

	resp = &types.FileUploadChunkCompleteReply{
		Identity: rp.Identity,
	}
	return
}
