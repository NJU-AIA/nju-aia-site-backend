package article

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	db.AutoMigrate(&Article{})
	return &Repository{db: db}
}

func (r *Repository) Create(a *Article) error {
	return r.db.Create(a).Error
}

// ... FindByID 和 List 保持不变 ...

func (r *Repository) FindByID(id string) (*Article, error) {
	var a Article
	// 注意：当 ID 是字符串时，建议显式写出查询条件 "id = ?"
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (r *Repository) List() ([]Article, error) {
	var articles []Article
	err := r.db.Order("created_at desc").Find(&articles).Error
	return articles, err
}