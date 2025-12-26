package models

import (
	"time"

	"gorm.io/gorm"
)

type StorageOrder struct {
	ID                    int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity              string         `gorm:"column:identity"`
	UserIdentity          string         `gorm:"column:user_identity"`
	StripeSessionID       string         `gorm:"column:stripe_session_id"`
	StripePaymentIntentID string         `gorm:"column:stripe_payment_intent_id"`
	StorageAmount         int64          `gorm:"column:storage_amount"` // Storage capacity in bytes
	PriceAmount           int64          `gorm:"column:price_amount"`   // Price in cents
	Currency              string         `gorm:"column:currency;default:usd"`
	Status                string         `gorm:"column:status;default:pending"` // pending, paid, failed, refunded
	CreatedAt             time.Time      `gorm:"column:created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (StorageOrder) TableName() string {
	return "storage_orders"
}
