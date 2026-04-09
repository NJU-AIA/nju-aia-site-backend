package article

import "time"

type Article struct {
	// ID
	ID string `gorm:"primaryKey;size:8" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Cover 文章封面图 URL
	Cover string `gorm:"size:512" json:"cover" example:"https://example.com/cover.jpg"`

	// Title 文章标题
	Title string `gorm:"size:255;not null" json:"title" example:"关于后端架构的深度思考"`

	// Content 文章正文
	Content string `gorm:"type:longtext;not null" json:"content" example:"# 第一章\n内容如下..."`

	// Author 作者姓名或昵称
	Author string `gorm:"size:100;not null" json:"author" example:"负责人"`

	// Category 文章分类
	Category string `gorm:"size:100;not null" json:"category" example:"技术分享"`

	// DefaultMode 展示模式：article, slide, homework
	DefaultMode string `gorm:"size:20;not null;default:'article'" json:"defaultMode" example:"article"`

	// Date 手动设置的展示日期
	Date string `gorm:"size:10;not null" json:"date" example:"2026-04-09"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateRequest 定义了发布或编辑文章时接收的参数格式
type CreateRequest struct {
	Title       string `json:"title" binding:"required" example:"标题"`
	Content     string `json:"content" binding:"required" example:"# 正文"`
	Author      string `json:"author" binding:"required" example:"作者"`
	Category    string `json:"category" binding:"required" example:"分类"`
	DefaultMode string `json:"defaultMode" binding:"required,oneof=article slide homework" example:"article"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02" example:"2026-04-09"`
	Cover       string `json:"cover"`
}
