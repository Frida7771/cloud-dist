package logic

import (
	"context"

	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
)

type StorageOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageOrderListLogic {
	return &StorageOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StorageOrderListLogic) StorageOrderList(req *types.StorageOrderListRequest, userIdentity string) (resp *types.StorageOrderListReply, err error) {
	resp = &types.StorageOrderListReply{
		List: []*types.StorageOrderItem{},
	}

	query := l.svcCtx.DB.WithContext(l.ctx).
		Model(&models.StorageOrder{}).
		Where("user_identity = ?", userIdentity)

	// Filter by status if provided
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Order by created_at descending
	query = query.Order("created_at DESC")

	var orders []models.StorageOrder
	if err = query.Find(&orders).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	for _, order := range orders {
		resp.List = append(resp.List, &types.StorageOrderItem{
			Identity:      order.Identity,
			StorageAmount: order.StorageAmount,
			PriceAmount:   order.PriceAmount,
			Currency:      order.Currency,
			Status:        order.Status,
			CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     order.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return
}
