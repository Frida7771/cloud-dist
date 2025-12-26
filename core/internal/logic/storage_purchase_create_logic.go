package logic

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"gorm.io/gorm"
)

type StoragePurchaseCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStoragePurchaseCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoragePurchaseCreateLogic {
	return &StoragePurchaseCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StoragePurchaseCreateLogic) StoragePurchaseCreate(req *types.StoragePurchaseCreateRequest, userIdentity string) (resp *types.StoragePurchaseCreateReply, err error) {
	// Validate storage amount
	if req.StorageAmount <= 0 {
		return nil, errors.New("storage amount must be greater than 0")
	}

	// Set Stripe API key
	if define.StripeSecretKey == "" {
		return nil, errors.New("Stripe secret key not configured")
	}
	stripe.Key = define.StripeSecretKey

	// Get price for the storage amount
	priceAmount := define.GetStoragePrice(req.StorageAmount)
	if priceAmount <= 0 {
		return nil, errors.New("invalid storage amount or pricing not configured")
	}

	// Set default currency
	currency := req.Currency
	if currency == "" {
		currency = "usd"
	}

	// Get user email for Stripe checkout
	ub := new(models.UserBasic)
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Select("email").Where("identity = ?", userIdentity).First(ub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Create order record
	order := &models.StorageOrder{
		Identity:      helper.UUID(),
		UserIdentity:  userIdentity,
		StorageAmount: req.StorageAmount,
		PriceAmount:   priceAmount,
		Currency:      currency,
		Status:        "pending",
	}
	if err = l.svcCtx.DB.WithContext(l.ctx).Create(order).Error; err != nil {
		log.Printf("[StoragePurchaseCreate] Failed to create order: %v", err)
		return nil, errors.New("failed to create order")
	}

	// Create Stripe Checkout Session
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(fmt.Sprintf("Storage: %s", formatStorageSize(req.StorageAmount))),
						Description: stripe.String(fmt.Sprintf("Additional storage capacity: %s", formatStorageSize(req.StorageAmount))),
					},
					UnitAmount: stripe.Int64(priceAmount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:3000/profile?payment=success"),
		CancelURL:  stripe.String("http://localhost:3000/profile?payment=cancel"),
		Metadata: map[string]string{
			"order_identity": order.Identity,
			"user_identity":  userIdentity,
			"storage_amount": fmt.Sprintf("%d", req.StorageAmount),
		},
	}

	// Add customer email if available
	if ub.Email != "" {
		params.CustomerEmail = stripe.String(ub.Email)
	}

	session, err := session.New(params)
	if err != nil {
		log.Printf("[StoragePurchaseCreate] Failed to create Stripe session: %v", err)
		// Update order status to failed
		l.svcCtx.DB.WithContext(l.ctx).Model(order).Update("status", "failed")
		return nil, errors.New("failed to create payment session")
	}

	// Update order with Stripe session ID
	order.StripeSessionID = session.ID
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(order).Update("stripe_session_id", session.ID).Error; err != nil {
		log.Printf("[StoragePurchaseCreate] Failed to update order with session ID: %v", err)
	}

	resp = &types.StoragePurchaseCreateReply{
		SessionID: session.ID,
		URL:       session.URL,
	}

	// Store session ID in response for frontend to use in success callback
	log.Printf("[StoragePurchaseCreate] Created session: %s, URL: %s", session.ID, session.URL)
	return
}

// formatStorageSize formats bytes into human-readable format
func formatStorageSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
