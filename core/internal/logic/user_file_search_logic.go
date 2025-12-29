package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"cloud-dist/core/define"
	"cloud-dist/core/internal/types"
	"cloud-dist/core/models"
	"cloud-dist/core/svc"

	"gorm.io/gorm"
)

type UserFileSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileSearchLogic {
	return &UserFileSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileSearchLogic) UserFileSearch(req *types.UserFileSearchRequest, userIdentity string) (resp *types.UserFileSearchReply, err error) {
	resp = new(types.UserFileSearchReply)
	resp.List = make([]*types.UserFileSearchItem, 0)

	// Validate keyword
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return resp, nil
	}

	size := req.Size
	if size == 0 {
		size = define.PageSize
	}
	page := req.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * size

	// Build base query for user's files
	query := l.svcCtx.DB.WithContext(l.ctx).Table("user_repository").
		Where("user_repository.user_identity = ?", userIdentity).
		Where("user_repository.deleted_at IS NULL").
		Where("user_repository.repository_identity != '' AND user_repository.repository_identity IS NOT NULL").
		Joins("LEFT JOIN repository_pool ON user_repository.repository_identity = repository_pool.identity")

	// Search by file name (case-insensitive, partial match)
	keywordPattern := "%" + keyword + "%"
	query = query.Where("user_repository.name LIKE ?", keywordPattern)

	// Filter by file type if provided
	if req.FileType != "" {
		fileType := strings.TrimSpace(req.FileType)
		if !strings.HasPrefix(fileType, ".") {
			fileType = "." + fileType
		}
		query = query.Where("user_repository.ext = ?", fileType)
	}

	// Select fields
	query = query.Select("user_repository.id, user_repository.identity, user_repository.repository_identity, " +
		"user_repository.ext, user_repository.name, repository_pool.path, repository_pool.size, user_repository.created_at, " +
		"user_repository.parent_id")

	// Get total count
	var count int64
	countQuery := query
	if err = countQuery.Count(&count).Error; err != nil {
		return nil, err
	}

	// Get paginated results
	var results []struct {
		ID                 int64
		Identity           string
		RepositoryIdentity string
		Ext                string
		Name               string
		Path               string
		Size               int64
		CreatedAt          time.Time
		ParentId           int64
	}

	if err = query.
		Order("user_repository.created_at DESC").
		Limit(size).
		Offset(offset).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// Build parent path for each file
	parentPathMap := make(map[int64]string)
	for _, result := range results {
		if result.ParentId > 0 {
			path := l.getParentPath(result.ParentId, userIdentity, parentPathMap)
			parentPathMap[result.ParentId] = path
		}
	}

	// Convert to response format
	for _, r := range results {
		item := &types.UserFileSearchItem{
			ID:                 r.ID,
			Identity:           r.Identity,
			RepositoryIdentity: r.RepositoryIdentity,
			Ext:                r.Ext,
			Name:               r.Name,
			Size:               r.Size,
			CreatedAt:          r.CreatedAt.Format(define.Datetime),
			ParentId:           r.ParentId,
		}

		if r.RepositoryIdentity != "" {
			item.Path = "/file/download?identity=" + r.RepositoryIdentity
		}

		if r.ParentId > 0 {
			if path, ok := parentPathMap[r.ParentId]; ok {
				item.ParentPath = path
			} else {
				item.ParentPath = l.getParentPath(r.ParentId, userIdentity, parentPathMap)
			}
			
			// Get parent folder identity
			var parentFolder models.UserRepository
			err := l.svcCtx.DB.WithContext(l.ctx).
				Select("identity").
				Where("id = ? AND user_identity = ?", r.ParentId, userIdentity).
				First(&parentFolder).Error
			if err == nil {
				item.ParentIdentity = parentFolder.Identity
			}
		} else {
			item.ParentPath = "Root"
			item.ParentIdentity = ""
		}

		resp.List = append(resp.List, item)
	}

	resp.Count = count
	return
}

func (l *UserFileSearchLogic) getParentPath(parentId int64, userIdentity string, cache map[int64]string) string {
	if path, ok := cache[parentId]; ok {
		return path
	}

	var parent models.UserRepository
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND user_identity = ?", parentId, userIdentity).
		First(&parent).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "Root"
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "Root"
	}

	if parent.ParentId == 0 {
		path := "Root/" + parent.Name
		cache[parentId] = path
		return path
	}

	parentPath := l.getParentPath(parent.ParentId, userIdentity, cache)
	path := parentPath + "/" + parent.Name
	cache[parentId] = path
	return path
}
