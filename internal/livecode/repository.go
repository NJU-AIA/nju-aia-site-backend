package livecode

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	_ = db.AutoMigrate(&Document{}, &Block{})
	return &Repository{db: db}
}

func (r *Repository) CreateDocument(doc *Document) error {
	return r.db.Create(doc).Error
}

func (r *Repository) FindDocumentByID(id string) (*Document, error) {
	var doc Document
	if err := r.db.First(&doc, "id = ?", id).Error; err != nil {
		return &doc, err
	}

	blocks, err := r.findBlocksByIDs(doc.BlockIDs)
	if err != nil {
		return nil, err
	}
	doc.Blocks = blocks
	return &doc, nil
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

func (r *Repository) AddBlock(documentID string, req BlockRequest) (*Block, error) {
	var created Block
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var doc Document
		if err := tx.First(&doc, "id = ?", documentID).Error; err != nil {
			return err
		}

		created = Block{
			ID:       blockID(),
			Type:     req.Type,
			Content:  req.Content,
			Language: normalizeLanguage(req.Type, req.Language),
		}

		if err := tx.Create(&created).Error; err != nil {
			return err
		}

		doc.BlockIDs = append(doc.BlockIDs, created.ID)
		return tx.Save(&doc).Error
	})
	return &created, err
}

func (r *Repository) UpdateBlock(documentID, blockID string, req BlockRequest) (*Block, error) {
	var updated Block
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var doc Document
		if err := tx.First(&doc, "id = ?", documentID).Error; err != nil {
			return err
		}
		if !containsString(doc.BlockIDs, blockID) {
			return gorm.ErrRecordNotFound
		}

		if err := tx.First(&updated, "id = ?", blockID).Error; err != nil {
			return err
		}

		updated.Type = req.Type
		updated.Content = req.Content
		updated.Language = normalizeLanguage(req.Type, req.Language)

		return tx.Save(&updated).Error
	})
	return &updated, err
}

func (r *Repository) UpdateBlockIDs(documentID string, blockIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var doc Document
		if err := tx.First(&doc, "id = ?", documentID).Error; err != nil {
			return err
		}

		if _, err := r.findBlocksByIDs(blockIDs); err != nil {
			return err
		}

		doc.BlockIDs = append([]string(nil), blockIDs...)
		return tx.Save(&doc).Error
	})
}

func (r *Repository) DeleteBlock(documentID, blockID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var doc Document
		if err := tx.First(&doc, "id = ?", documentID).Error; err != nil {
			return err
		}
		if !containsString(doc.BlockIDs, blockID) {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Delete(&Block{}, "id = ?", blockID).Error; err != nil {
			return err
		}

		doc.BlockIDs = removeString(doc.BlockIDs, blockID)
		return tx.Save(&doc).Error
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

func (r *Repository) findBlocksByIDs(ids []string) ([]Block, error) {
	if len(ids) == 0 {
		return []Block{}, nil
	}

	var blocks []Block
	if err := r.db.Where("id IN ?", ids).Find(&blocks).Error; err != nil {
		return nil, err
	}

	blockMap := make(map[string]Block, len(blocks))
	for _, block := range blocks {
		blockMap[block.ID] = block
	}

	ordered := make([]Block, 0, len(ids))
	for _, id := range ids {
		block, ok := blockMap[id]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		ordered = append(ordered, block)
	}

	return ordered, nil
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func removeString(values []string, target string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != target {
			result = append(result, value)
		}
	}
	return result
}
