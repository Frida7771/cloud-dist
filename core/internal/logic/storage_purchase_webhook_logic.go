package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"cloud-dist/core/define"
	"cloud-dist/core/svc"
	"cloud-dist/core/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
	"gorm.io/gorm"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type StoragePurchaseWebhookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStoragePurchaseWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoragePurchaseWebhookLogic {
	return &StoragePurchaseWebhookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StoragePurchaseWebhookLogic) ProcessWebhook(payload []byte, signature string) error {
	// Set Stripe API key
	if define.StripeSecretKey == "" {
		return errors.New("Stripe secret key not configured")
	}
	stripe.Key = define.StripeSecretKey

	// Verify webhook secret is configured
	if define.StripeWebhookSecret == "" {
		log.Printf("[StoragePurchaseWebhook] ERROR: Stripe webhook secret not configured")
		return errors.New("Stripe webhook secret not configured")
	}

	log.Printf("[StoragePurchaseWebhook] Verifying signature with secret: %s...%s (len: %d)",
		define.StripeWebhookSecret[:10],
		define.StripeWebhookSecret[len(define.StripeWebhookSecret)-10:],
		len(define.StripeWebhookSecret))

	sigPreview := signature
	if len(signature) > 100 {
		sigPreview = signature[:100] + "..."
	}
	log.Printf("[StoragePurchaseWebhook] Payload length: %d, Signature preview: %s", len(payload), sigPreview)

	// Use ConstructEventWithOptions to ignore API version mismatch
	// Stripe CLI may send events with newer API versions than the SDK expects
	event, err := webhook.ConstructEventWithOptions(payload, signature, define.StripeWebhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		log.Printf("[StoragePurchaseWebhook] ERROR: Failed to verify webhook signature: %v", err)
		log.Printf("[StoragePurchaseWebhook] Error details: %+v", err)
		if len(payload) > 0 {
			previewLen := min(200, len(payload))
			log.Printf("[StoragePurchaseWebhook] Payload (first %d bytes): %s", previewLen, string(payload[:previewLen]))
		}
		log.Printf("[StoragePurchaseWebhook] WebhookSecret length: %d", len(define.StripeWebhookSecret))
		log.Printf("[StoragePurchaseWebhook] Signature length: %d", len(signature))
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	log.Printf("[StoragePurchaseWebhook] Received event: %s (ID: %s)", event.Type, event.ID)

	// Handle the event
	switch event.Type {
	case "checkout.session.completed":
		return l.handleCheckoutSessionCompleted(event)
	case "payment_intent.succeeded":
		return l.handlePaymentIntentSucceeded(event)
	case "payment_intent.payment_failed":
		return l.handlePaymentIntentFailed(event)
	default:
		log.Printf("[StoragePurchaseWebhook] Unhandled event type: %s", event.Type)
	}

	return nil
}

func (l *StoragePurchaseWebhookLogic) handleCheckoutSessionCompleted(event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		log.Printf("[StoragePurchaseWebhook] Failed to unmarshal checkout session: %v", err)
		return err
	}

	log.Printf("[StoragePurchaseWebhook] Checkout session completed: %s, PaymentStatus: %s, Mode: %s",
		session.ID, session.PaymentStatus, session.Mode)

	// Get order by Stripe session ID
	order := new(models.StorageOrder)
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("stripe_session_id = ?", session.ID).First(order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[StoragePurchaseWebhook] Order not found for session: %s", session.ID)
			return nil // Not an error, just log it
		}
		return err
	}

	log.Printf("[StoragePurchaseWebhook] Found order: %s, current status: %s", order.Identity, order.Status)

	// If order is already paid, skip
	if order.Status == "paid" {
		log.Printf("[StoragePurchaseWebhook] Order %s is already paid, skipping", order.Identity)
		return nil
	}

	// For checkout.session.completed event, the session is completed which means payment was successful
	// We can update the order status directly
	// Note: PaymentStatus might be "paid" or "unpaid" depending on timing, but session.completed means payment succeeded
	log.Printf("[StoragePurchaseWebhook] Updating order %s to paid (checkout.session.completed)", order.Identity)
	order.Status = "paid"
	if session.PaymentIntent != nil {
		order.StripePaymentIntentID = session.PaymentIntent.ID
		log.Printf("[StoragePurchaseWebhook] Payment Intent ID: %s", session.PaymentIntent.ID)
	} else {
		// If PaymentIntent is not in the event, we might need to expand it
		// But for now, we'll update anyway and payment_intent.succeeded will also handle it
		log.Printf("[StoragePurchaseWebhook] PaymentIntent not found in session, will be set by payment_intent.succeeded event")
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Save(order).Error; err != nil {
		log.Printf("[StoragePurchaseWebhook] Failed to update order status: %v", err)
		return err
	}

	// Update user's total volume
	if err := l.updateUserStorage(order.UserIdentity, order.StorageAmount); err != nil {
		log.Printf("[StoragePurchaseWebhook] Failed to update user storage: %v", err)
		return err
	}

	log.Printf("[StoragePurchaseWebhook] Order %s marked as paid, user storage updated", order.Identity)
	return nil
}

func (l *StoragePurchaseWebhookLogic) handlePaymentIntentSucceeded(event stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		log.Printf("[StoragePurchaseWebhook] Failed to unmarshal payment intent: %v", err)
		return err
	}

	log.Printf("[StoragePurchaseWebhook] Payment intent succeeded: %s", paymentIntent.ID)

	// Find order by payment intent ID (if already set from checkout.session.completed)
	order := new(models.StorageOrder)
	found := false

	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("stripe_payment_intent_id = ?", paymentIntent.ID).First(order).Error; err == nil {
		found = true
		log.Printf("[StoragePurchaseWebhook] Found order by payment intent ID: %s", order.Identity)
	}

	// If not found, try to find by metadata (order_identity)
	if !found && paymentIntent.Metadata != nil {
		if orderIdentity, ok := paymentIntent.Metadata["order_identity"]; ok {
			if err := l.svcCtx.DB.WithContext(l.ctx).
				Where("identity = ?", orderIdentity).First(order).Error; err == nil {
				found = true
				log.Printf("[StoragePurchaseWebhook] Found order by metadata order_identity: %s", order.Identity)
			}
		}
	}

	// If still not found, search through recent pending orders and match by session
	if !found {
		log.Printf("[StoragePurchaseWebhook] Searching for order by matching payment intent in recent sessions...")
		var pendingOrders []models.StorageOrder
		if err := l.svcCtx.DB.WithContext(l.ctx).
			Where("status = ?", "pending").
			Where("stripe_session_id != ?", "").
			Order("created_at DESC").
			Limit(20).
			Find(&pendingOrders).Error; err == nil {
			// Query Stripe API to get session details and match payment intent
			for _, pendingOrder := range pendingOrders {
				if pendingOrder.StripeSessionID != "" {
					// Query Stripe to get session details
					stripe.Key = define.StripeSecretKey
					sessionParams := &stripe.CheckoutSessionParams{}
					sessionParams.AddExpand("payment_intent")
					if stripeSession, err := session.Get(pendingOrder.StripeSessionID, sessionParams); err == nil {
						if stripeSession.PaymentIntent != nil && stripeSession.PaymentIntent.ID == paymentIntent.ID {
							order = &pendingOrder
							found = true
							log.Printf("[StoragePurchaseWebhook] Found order by matching payment intent in session: %s", order.Identity)
							break
						}
					}
				}
			}
		}
	}

	if !found {
		log.Printf("[StoragePurchaseWebhook] Order not found for payment intent: %s", paymentIntent.ID)
		return nil
	}

	// Update order status if still pending
	if order.Status == "pending" {
		log.Printf("[StoragePurchaseWebhook] Updating order %s from pending to paid", order.Identity)
		order.Status = "paid"
		order.StripePaymentIntentID = paymentIntent.ID
		if err := l.svcCtx.DB.WithContext(l.ctx).Save(order).Error; err != nil {
			log.Printf("[StoragePurchaseWebhook] Failed to update order status: %v", err)
			return err
		}

		// Update user's total volume
		if err := l.updateUserStorage(order.UserIdentity, order.StorageAmount); err != nil {
			log.Printf("[StoragePurchaseWebhook] Failed to update user storage: %v", err)
			return err
		}

		log.Printf("[StoragePurchaseWebhook] Order %s marked as paid, user storage updated", order.Identity)
	} else if order.Status == "paid" {
		log.Printf("[StoragePurchaseWebhook] Order %s is already paid, skipping", order.Identity)
	}

	return nil
}

func (l *StoragePurchaseWebhookLogic) handlePaymentIntentFailed(event stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		log.Printf("[StoragePurchaseWebhook] Failed to unmarshal payment intent: %v", err)
		return err
	}

	log.Printf("[StoragePurchaseWebhook] Payment intent failed: %s", paymentIntent.ID)

	// Find and update order status
	order := new(models.StorageOrder)
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("stripe_payment_intent_id = ?", paymentIntent.ID).First(order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[StoragePurchaseWebhook] Order not found for payment intent: %s", paymentIntent.ID)
			return nil
		}
		return err
	}

	if order.Status == "pending" {
		order.Status = "failed"
		if err := l.svcCtx.DB.WithContext(l.ctx).Save(order).Error; err != nil {
			log.Printf("[StoragePurchaseWebhook] Failed to update order status: %v", err)
			return err
		}
		log.Printf("[StoragePurchaseWebhook] Order %s marked as failed", order.Identity)
	}

	return nil
}

func (l *StoragePurchaseWebhookLogic) updateUserStorage(userIdentity string, additionalStorage int64) error {
	// Update user's total volume by adding the purchased storage
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&models.UserBasic{}).
		Where("identity = ?", userIdentity).
		UpdateColumn("total_volume", gorm.Expr("total_volume + ?", additionalStorage)).Error; err != nil {
		return fmt.Errorf("failed to update user storage: %w", err)
	}

	log.Printf("[StoragePurchaseWebhook] User %s storage increased by %d bytes", userIdentity, additionalStorage)
	return nil
}
