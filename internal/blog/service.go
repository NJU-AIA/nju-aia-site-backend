package blog

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveMarkdown(title, author, content string) (uint, error) {
	b := &Blog{
		Title:   title,
		Author:  author,
		Content: content,
	}
	err := s.repo.Create(b)
	return b.ID, err
}

func (s *Service) GetBlog(id uint) (*Blog, error) {
	return s.repo.FindByID(id)
}