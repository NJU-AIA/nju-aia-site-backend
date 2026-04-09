package livecode

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	_ = db.AutoMigrate(&Document{}, &Block{})
	return &Repository{db: db}
}

func (r *Repository) CreateDocumentWithBlocks(doc *Document, blocks []Block) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(doc).Error; err != nil {
			return err
		}
		if len(blocks) == 0 {
			return nil
		}
		return tx.Create(&blocks).Error
	})
}

func (r *Repository) FindDocumentByID(id string) (*Document, error) {
	var doc Document
	err := r.db.Preload("Blocks", func(db *gorm.DB) *gorm.DB {
		return db.Order("block_order asc")
	}).First(&doc, "id = ?", id).Error
	return &doc, err
}

func (r *Repository) ListDocuments() ([]ListItem, error) {
	var items []ListItem
	err := r.db.Model(&Document{}).
		Select("id", "slug", "published_at").
		Order("published_at desc").
		Order("id desc").
		Find(&items).Error
	return items, err
}

func (r *Repository) UpdateDocumentWithBlocks(doc *Document, blocks []Block) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Document{}).
			Where("id = ?", doc.ID).
			Updates(map[string]any{
				"slug":         doc.Slug,
				"published_at": doc.PublishedAt,
			}).Error; err != nil {
			return err
		}

		if err := tx.Where("document_id = ?", doc.ID).Delete(&Block{}).Error; err != nil {
			return err
		}

		if len(blocks) == 0 {
			return nil
		}

		return tx.Create(&blocks).Error
	})
}

func (r *Repository) DeleteDocumentByID(id string) error {
	return r.db.Delete(&Document{}, "id = ?", id).Error
}

func (r *Repository) ExistsSlug(slug string, excludeID string) (bool, error) {
	var count int64
	db := r.db.Model(&Document{}).Where("slug = ?", slug)
	if excludeID != "" {
		db = db.Where("id <> ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
