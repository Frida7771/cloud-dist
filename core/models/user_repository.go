package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	ID                 int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity           string         `gorm:"column:identity"`
	UserIdentity       string         `gorm:"column:user_identity"`
	ParentId           int64          `gorm:"column:parent_id"`
	RepositoryIdentity string         `gorm:"column:repository_identity"`
	Ext                string         `gorm:"column:ext"`
	Name               string         `gorm:"column:name"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (UserRepository) TableName() string {
	return "user_repository"
}
