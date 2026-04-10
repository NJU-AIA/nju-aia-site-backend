package article

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

// Service 处理文章相关的业务逻辑
type Service struct {
	repo *Repository
}

// NewService 初始化业务逻辑层
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateArticle 创建文章的业务逻辑
func (s *Service) CreateArticle(req CreateRequest) (string, error) {
	fullID := strings.ReplaceAll(uuid.New().String(), "-", "")
	shortID := fullID[:8]

	newArt := &Article{
		ID:          shortID,
		Title:       req.Title,
		Category:    req.Category,
		Author:      req.Author,
		DefaultMode: req.DefaultMode,
		Published:   req.Published,
		Date:        req.Date,
		Content:     req.Content,
		Cover:       req.Cover,
	}

	err := s.repo.Create(newArt)
	return newArt.ID, err
}

// GetArticle 获取单篇文章详情
func (s *Service) GetArticle(id string, publishedOnly bool) (*Article, error) {
	return s.repo.FindByID(id, publishedOnly)
}

// UpdateArticle 更新文章的业务逻辑
func (s *Service) UpdateArticle(id string, req CreateRequest) error {
	// 1. 检查文章是否存在
	art, err := s.repo.FindByID(id, false)
	if err != nil {
		return err
	}

	// 2. 覆盖数据
	art.Title = req.Title
	art.Category = req.Category
	art.Author = req.Author
	art.DefaultMode = req.DefaultMode
	art.Published = req.Published
	art.Date = req.Date
	art.Content = req.Content
	art.Cover = req.Cover
	// 3. 调用仓库执行保存
	return s.repo.Update(art)
}

// DeleteArticle 删除文章的业务逻辑
func (s *Service) DeleteArticle(id string) error {
	_, err := s.repo.FindByID(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return s.repo.DeleteByID(id)
}
