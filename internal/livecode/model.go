package livecode

import "time"

type Document struct {
	ID          string    `gorm:"primaryKey;size:8" json:"id" example:"lc001234"`
	Slug        string    `gorm:"size:100;uniqueIndex;not null" json:"slug" example:"python-basic-demo"`
	PublishedAt string    `gorm:"size:10;not null" json:"publishedAt" example:"2026-04-09"`
	BlockIDs    []string  `gorm:"serializer:json;not null" json:"blockIds"`
	Blocks      []Block   `gorm:"-" json:"blocks,omitempty"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Block struct {
	ID        string    `gorm:"primaryKey;size:12" json:"id" example:"blk000000001"`
	Type      string    `gorm:"size:20;not null" json:"type" example:"code"`
	Content   string    `gorm:"type:longtext;not null" json:"content" example:"print('hello world')"`
	Language  string    `gorm:"size:50" json:"language,omitempty" example:"python"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UpsertRequest struct {
	Slug        string `json:"slug" binding:"required" example:"python-basic-demo"`
	PublishedAt string `json:"publishedAt" binding:"required,datetime=2006-01-02" example:"2026-04-09"`
}

type BlockRequest struct {
	Type     string `json:"type" binding:"required,oneof=markdown code" example:"code"`
	Content  string `json:"content" binding:"required" example:"print('hello world')"`
	Language string `json:"language,omitempty" example:"python"`
}

type BlockIDsRequest struct {
	BlockIDs []string `json:"blockIds" binding:"required,min=0,dive,required" example:"[\"blk1\",\"blk2\"]"`
}

type BlockResponse struct {
	ID       string `json:"id" example:"blk000000001"`
	Type     string `json:"type" example:"code"`
	Content  string `json:"content" example:"print('hello world')"`
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
