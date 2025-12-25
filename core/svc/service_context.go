package svc

import (
	"context"
	"fmt"

	"cloud-disk/core/define"
	"cloud-disk/core/internal/middleware"
	"cloud-disk/core/models"
	appcfg "cloud-disk/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config appcfg.Config
	DB     *gorm.DB
	RDB    *redis.Client
	Auth   gin.HandlerFunc
}

func NewServiceContext(c appcfg.Config) (*ServiceContext, error) {
	db, err := models.Init(c.Mysql.DataSource)
	if err != nil {
		return nil, fmt.Errorf("init mysql: %w", err)
	}
	rdb := models.InitRedis(c.Redis.Addr)

	// Initialize S3 configuration from config file (environment variables take precedence)
	define.InitS3Config(
		c.S3.AccessKeyID,
		c.S3.SecretAccessKey,
		c.S3.Bucket,
		c.S3.Region,
		c.S3.Endpoint,
	)

	return &ServiceContext{
		Config: c,
		DB:     db,
		RDB:    rdb,
		Auth:   middleware.NewAuthMiddleware().Handle,
	}, nil
}

func (s *ServiceContext) Close(ctx context.Context) error {
	if s.DB != nil {
		sqlDB, err := s.DB.DB()
		if err == nil {
			if err = sqlDB.Close(); err != nil {
				return fmt.Errorf("close mysql: %w", err)
			}
		}
	}
	if s.RDB != nil {
		if err := s.RDB.Close(); err != nil {
			return fmt.Errorf("close redis: %w", err)
		}
	}
	return nil
}
