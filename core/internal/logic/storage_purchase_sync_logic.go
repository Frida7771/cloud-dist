package logic

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud-dist/core/define"
	"cloud-dist/core/svc"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"gorm.io/gorm"
)

type StoragePurchaseSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStoragePurchaseSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoragePurchaseSyncLogic {
	return &StoragePurchaseSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StoragePurchaseSyncLogic) StoragePurchaseSync(req *types.StoragePurchaseSyncRequest) (resp *types.StoragePurchaseSyncReply, err error) {
	resp = &types.StoragePurchaseSyncReply{
		Status:  "pending",
		Message: "Order not found or payment not completed",
	}

	// Set Stripe API key
	if define.StripeSecretKey == "" {
		return nil, errors.New("Stripe secret key not configured")
	}
	stripe.Key = define.StripeSecretKey

	// Get order by Stripe session ID
	order := new(models.StorageOrder)
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Where("stripe_session_id = ?", req.SessionID).First(order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// If already paid, return immediately
	if order.Status == "paid" {
		resp.Status = "paid"
		resp.StorageAmount = order.StorageAmount
		resp.Message = "Order already paid"
		return resp, nil
	}

	// Query Stripe to check payment status
	sessionParams := &stripe.CheckoutSessionParams{}
	sessionParams.AddExpand("payment_intent")
	stripeSession, err := session.Get(req.SessionID, sessionParams)
	if err != nil {
		log.Printf("[StoragePurchaseSync] Failed to get Stripe session: %v", err)
		return nil, fmt.Errorf("failed to verify payment: %w", err)
	}

	// Check payment status
	if stripeSession.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		// Payment is successful, update order status
		order.Status = "paid"
		if stripeSession.PaymentIntent != nil {
			order.StripePaymentIntentID = stripeSession.PaymentIntent.ID
		}
		if err = l.svcCtx.DB.WithContext(l.ctx).Save(order).Error; err != nil {
			log.Printf("[StoragePurchaseSync] Failed to update order status: %v", err)
			return nil, errors.New("failed to update order status")
		}

		// Update user's total volume
		if err = l.updateUserStorage(order.UserIdentity, order.StorageAmount); err != nil {
			log.Printf("[StoragePurchaseSync] Failed to update user storage: %v", err)
			// Don't fail the request, just log the error
			// The webhook will handle it later
		}

		resp.Status = "paid"
		resp.StorageAmount = order.StorageAmount
		resp.Message = "Payment confirmed, storage capacity increased"
		log.Printf("[StoragePurchaseSync] Order %s synced and marked as paid", order.Identity)
	} else if stripeSession.PaymentStatus == stripe.CheckoutSessionPaymentStatusUnpaid {
		resp.Status = "pending"
		resp.Message = "Payment is still pending"
	} else {
		// Payment failed or cancelled
		order.Status = "failed"
		if err = l.svcCtx.DB.WithContext(l.ctx).Save(order).Error; err != nil {
			log.Printf("[StoragePurchaseSync] Failed to update order status: %v", err)
		}
		resp.Status = "failed"
		resp.Message = "Payment failed or was cancelled"
	}

	return resp, nil
}

func (l *StoragePurchaseSyncLogic) updateUserStorage(userIdentity string, additionalStorage int64) error {
	// Update user's total volume by adding the purchased storage
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		UpdateColumn("total_volume", gorm.Expr("total_volume + ?", additionalStorage)).Error; err != nil {
		return fmt.Errorf("failed to update user storage: %w", err)
	}

	log.Printf("[StoragePurchaseSync] User %s storage increased by %d bytes", userIdentity, additionalStorage)
	return nil
}


