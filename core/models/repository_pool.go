package models

import (
	"time"

	"gorm.io/gorm"
)

type RepositoryPool struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity  string         `gorm:"column:identity"`
	Hash      string         `gorm:"column:hash"`
	Name      string         `gorm:"column:name"`
	Ext       string         `gorm:"column:ext"`
	Size      int64          `gorm:"column:size"`
	Path      string         `gorm:"column:path"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (RepositoryPool) TableName() string {
	return "repository_pool"
}
