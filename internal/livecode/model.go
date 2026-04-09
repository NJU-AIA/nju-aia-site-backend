package livecode

import "time"

type Document struct {
	ID          string    `gorm:"primaryKey;size:8" json:"id" example:"lc001234"`
	Slug        string    `gorm:"size:100;uniqueIndex;not null" json:"slug" example:"python-basic-demo"`
	PublishedAt string    `gorm:"size:10;not null" json:"publishedAt" example:"2026-04-09"`
	Blocks      []Block   `gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE" json:"blocks,omitempty"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Block struct {
	ID         string    `gorm:"primaryKey;size:12" json:"id" example:"blk000000001"`
	DocumentID string    `gorm:"size:8;index;not null" json:"-"`
	Type       string    `gorm:"size:20;not null" json:"type" example:"code"`
	Order      int       `gorm:"column:block_order;not null;index" json:"order" example:"1"`
	Content    string    `gorm:"type:longtext;not null" json:"content" example:"print('hello world')"`
	Language   string    `gorm:"size:50" json:"language,omitempty" example:"python"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type UpsertRequest struct {
	Slug        string         `json:"slug" binding:"required" example:"python-basic-demo"`
	PublishedAt string         `json:"publishedAt" binding:"required,datetime=2006-01-02" example:"2026-04-09"`
	Blocks      []BlockRequest `json:"blocks" binding:"required,min=1,dive"`
}

type BlockRequest struct {
	Type     string `json:"type" binding:"required,oneof=markdown code" example:"code"`
	Order    int    `json:"order" binding:"required,min=1" example:"1"`
	Content  string `json:"content" binding:"required" example:"print('hello world')"`
	Language string `json:"language,omitempty" example:"python"`
}

type ListItem struct {
	ID          string `json:"id" example:"lc001234"`
	Slug        string `json:"slug" example:"python-basic-demo"`
	PublishedAt string `json:"publishedAt" example:"2026-04-09"`
}

type ListResponse struct {
	Items []ListItem `json:"items"`
}
