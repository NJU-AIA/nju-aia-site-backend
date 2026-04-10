package article

import "gorm.io/gorm"

// Repository 封装对 Article 表的数据库操作
type Repository struct {
	db *gorm.DB
}

// NewRepository 初始化仓库并自动同步表结构
func NewRepository(db *gorm.DB) *Repository {
	_ = db.AutoMigrate(&Article{})
	return &Repository{db: db}
}

// Create 在数据库中插入一条新文章记录
func (r *Repository) Create(a *Article) error {
	return r.db.Create(a).Error
}

// FindByID 根据 UUID 查找特定文章详情
func (r *Repository) FindByID(id string, publishedOnly bool) (*Article, error) {
	var a Article
	db := r.db
	if publishedOnly {
		db = db.Where("published = ?", true)
	}
	err := db.First(&a, "id = ?", id).Error
	return &a, err
}

// List 获取所有文章的元数据列表
func (r *Repository) List(publishedOnly bool) ([]Article, error) {
	var articles []Article
	db := r.db.Model(&Article{})
	if publishedOnly {
		db = db.Where("published = ?", true)
	}

	err := db.Select("id", "title", "category", "author", "default_mode", "published", "date", "created_at", "updated_at", "cover").
		Order("date desc").       // 按手动设置的展示日期倒序排列
		Order("created_at desc"). // 同日期下按创建时间倒序排列
		Find(&articles).Error

	return articles, err
}

// Update 更新数据库中的文章记录
func (r *Repository) Update(a *Article) error {
	return r.db.Save(a).Error
}

// DeleteByID 删除指定的文章记录
func (r *Repository) DeleteByID(id string) error {
	return r.db.Delete(&Article{}, "id = ?", id).Error
}
