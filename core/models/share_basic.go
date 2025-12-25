package models

import (
	"time"

	"gorm.io/gorm"
)

type ShareBasic struct {
	ID                     int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity               string         `gorm:"column:identity"`
	UserIdentity           string         `gorm:"column:user_identity"`
	UserRepositoryIdentity string         `gorm:"column:user_repository_identity"`
	RepositoryIdentity     string         `gorm:"column:repository_identity"`
	ExpiredTime            int            `gorm:"column:expired_time"`
	ClickNum               int            `gorm:"column:click_num"`
	CreatedAt              time.Time      `gorm:"column:created_at"`
	UpdatedAt              time.Time      `gorm:"column:updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ShareBasic) TableName() string {
	return "share_basic"
}
