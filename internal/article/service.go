package article

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateArticle 创建新文章
func (s *Service) CreateArticle(req CreateRequest) (string, error) {
	newArt := &Article{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Category:    req.Category,
		Author:      req.Author,
		DefaultMode: req.DefaultMode,
		Content:     req.Content,
	}
	err := s.repo.Create(newArt)
	return newArt.ID, err
}

func (s *Service) GetArticle(id string) (*Article, error) {
	// 直接传递字符串 id，不要做类型转换
	return s.repo.FindByID(id)
}