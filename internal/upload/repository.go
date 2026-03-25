package upload

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	db.AutoMigrate(&FileRecord{})
	return &Repository{db: db}
}

func (r *Repository) Create(file *FileRecord) error {
	return r.db.Create(file).Error
}