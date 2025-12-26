package models

import (
	"time"

	"gorm.io/gorm"
)

// Friend represents a friendship between two users
type Friend struct {
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity       string         `gorm:"column:identity"`
	UserIdentity   string         `gorm:"column:user_identity"`         // The user who has this friend
	FriendIdentity string         `gorm:"column:friend_identity"`       // The friend's user identity
	Status         string         `gorm:"column:status;default:active"` // active, blocked
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Friend) TableName() string {
	return "friend"
}

// FriendRequest represents a friend request
type FriendRequest struct {
	ID               int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity         string         `gorm:"column:identity"`
	FromUserIdentity string         `gorm:"column:from_user_identity"`     // User who sent the request
	ToUserIdentity   string         `gorm:"column:to_user_identity"`       // User who received the request
	Status           string         `gorm:"column:status;default:pending"` // pending, accepted, rejected
	Message          string         `gorm:"column:message"`                // Optional message
	CreatedAt        time.Time      `gorm:"column:created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (FriendRequest) TableName() string {
	return "friend_request"
}

// FriendShare represents a file shared with a friend
type FriendShare struct {
	ID                     int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Identity               string         `gorm:"column:identity"`
	FromUserIdentity       string         `gorm:"column:from_user_identity"`       // User who shared the file
	ToUserIdentity         string         `gorm:"column:to_user_identity"`         // Friend who received the share
	RepositoryIdentity     string         `gorm:"column:repository_identity"`      // The shared file
	UserRepositoryIdentity string         `gorm:"column:user_repository_identity"` // User's file reference
	Message                string         `gorm:"column:message"`                  // Optional message
	IsRead                 bool           `gorm:"column:is_read;default:false"`    // Whether the friend has read it
	CreatedAt              time.Time      `gorm:"column:created_at"`
	UpdatedAt              time.Time      `gorm:"column:updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (FriendShare) TableName() string {
	return "friend_share"
}
