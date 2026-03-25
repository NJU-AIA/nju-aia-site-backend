package asset

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	db.AutoMigrate(&AssetRecord{})
	return &Repository{db: db}
}

func (r *Repository) Create(asset *AssetRecord) error {
	return r.db.Create(asset).Error
}