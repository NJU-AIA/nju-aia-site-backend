package blog

import "time"

// Blog 博客实体
type Blog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Title     string    `gorm:"size:255;not null" json:"title" example:"标题五个字"`
	Content   string    `gorm:"type:text;not null" json:"content" example:"# Markdown内容"`
	Author    string    `gorm:"size:100;default:'匿名'" json:"author" example:"张三"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}