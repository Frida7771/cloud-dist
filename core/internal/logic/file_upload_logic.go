package logic

import (
	"context"
	"log"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
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
	log.Printf("[FileUploadLogic] 创建文件记录: 文件名=%s, 扩展名=%s, 大小=%d, MD5=%s", req.Name, req.Ext, req.Size, req.Hash)

	rp := &models.RepositoryPool{
		Identity: helper.UUID(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(rp).Error; err != nil {
		log.Printf("[FileUploadLogic] 保存文件记录失败: %v", err)
		return nil, err
	}
	log.Printf("[FileUploadLogic] 文件记录保存成功: identity=%s", rp.Identity)

	resp = &types.FileUploadReply{
		Identity: rp.Identity,
		Ext:      rp.Ext,
		Name:     rp.Name,
	}
	return
}
