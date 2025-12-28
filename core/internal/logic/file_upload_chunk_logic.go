package logic

import (
	"context"

	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
)

type FileUploadChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkLogic {
	return &FileUploadChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkLogic) FileUploadChunk(req *types.FileUploadChunkRequest) (resp *types.FileUploadChunkReply, err error) {
	return
}
