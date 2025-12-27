package logic

import (
	"context"
	"errors"

	"cloud-disk/core/define"
	"cloud-disk/core/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"gorm.io/gorm"
)

type UserFileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListRequest, userIdentity string) (resp *types.UserFileListReply, err error) {
	resp = new(types.UserFileListReply)

	size := req.Size
	if size == 0 {
		size = define.PageSize
	}
	page := req.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * size

	var parentID int64
	if req.Identity != "" {
		ur := new(models.UserRepository)
		err = l.svcCtx.DB.WithContext(l.ctx).Select("id").Where("identity = ?", req.Identity).First(ur).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if err == nil {
			parentID = ur.ID
		}
	}

	files := make([]*types.UserFile, 0)
	// For files (with repository_identity), use deduplication to avoid showing duplicates
	// For folders (without repository_identity), show all folders
	// Use subquery to get only one record per repository_identity for files (deduplication)
	subquery := l.svcCtx.DB.WithContext(l.ctx).Table("user_repository").
		Where("user_repository.user_identity = ?", userIdentity).
		Where("user_repository.parent_id = ?", parentID).
		Where("user_repository.deleted_at IS NULL").
		Where("user_repository.repository_identity != '' AND user_repository.repository_identity IS NOT NULL").
		Select("MIN(user_repository.id) as id").
		Group("user_repository.repository_identity")

	// Query: show all folders OR deduplicated files
	query := l.svcCtx.DB.WithContext(l.ctx).Table("user_repository").
		Where("user_repository.user_identity = ?", userIdentity).
		Where("user_repository.parent_id = ?", parentID).
		Where("user_repository.deleted_at IS NULL").
		Where("(user_repository.repository_identity = '' OR user_repository.repository_identity IS NULL OR user_repository.id IN (?))", subquery).
		Joins("LEFT JOIN repository_pool ON user_repository.repository_identity = repository_pool.identity").
		Select("user_repository.id, user_repository.identity, user_repository.repository_identity, user_repository.ext, " +
			"user_repository.name, repository_pool.path, repository_pool.size")

	if err = query.
		Limit(size).
		Offset(offset).
		Scan(&files).Error; err != nil {
		return nil, err
	}

	var count int64
	if err = l.svcCtx.DB.WithContext(l.ctx).Model(&models.UserRepository{}).
		Where("parent_id = ? AND user_identity = ?", parentID, userIdentity).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return nil, err
	}

	// Convert repository identity to download endpoint URL for each file
	// This provides permanent download links that don't expire
	for _, file := range files {
		if file.RepositoryIdentity != "" {
			// Use permanent download endpoint URL
			// This URL will work as long as user has permission
			file.Path = "/file/download?identity=" + file.RepositoryIdentity
		}
	}

	resp.List = files
	resp.Count = count

	return
}
