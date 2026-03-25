package blog

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	db.AutoMigrate(&Blog{})
	return &Repository{db: db}
}

func (r *Repository) Create(b *Blog) error {
	return r.db.Create(b).Error
}

func (r *Repository) FindByID(id uint) (*Blog, error) {
	var b Blog
	err := r.db.First(&b, id).Error
	return &b, err
}

func (r *Repository) List() ([]Blog, error) {
	var blogs []Blog
	err := r.db.Order("created_at desc").Find(&blogs).Error
	return blogs, err
}