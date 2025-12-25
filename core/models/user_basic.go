package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity    string         `gorm:"column:identity"`
	Name        string         `gorm:"column:name"`
	Password    string         `gorm:"column:password"`
	Email       string         `gorm:"column:email"`
	NowVolume   int64          `gorm:"column:now_volume"`
	TotalVolume int64          `gorm:"column:total_volume"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (UserBasic) TableName() string {
	return "user_basic"
}
