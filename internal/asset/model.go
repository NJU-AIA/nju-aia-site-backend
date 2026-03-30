package asset

import (
	"errors"
	"time"
)

var ErrUnsafePath = errors.New("scope 或 name 含有非法路径字符")

// AssetRecord 对应数据库中的静态资源记录
type AssetRecord struct {
	// ID 仅用于数据库内部主键
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"-"`

	// Path 完整资源路径
	Path          string    `gorm:"size:512;uniqueIndex;not null" json:"path"`

	// URL 完整访问地址
	URL           string    `gorm:"size:512;not null" json:"url"`

	// MarkdownValue 推荐写入 Markdown 的值
	MarkdownValue string    `gorm:"size:512;not null" json:"markdownValue"`

	// Scope 第一层路径
	Scope         string    `gorm:"size:100;index;not null" json:"scope"`

	// Kind 资源归属类型: article, shared
	Kind          string    `gorm:"size:20;index;not null" json:"kind"`

	// Filename 完整文件名
	Filename      string    `gorm:"size:255;not null" json:"filename"`

	// Name 文件名主干
	Name          string    `gorm:"size:255;not null" json:"name"`

	// Ext 扩展名
	Ext           string    `gorm:"size:20;not null" json:"ext"`

	// ContentType 文件的 MIME 类型
	ContentType   string    `gorm:"size:100" json:"contentType"`

	// Size 文件大小
	Size          int64     `json:"size"`

	// SavedPath 存储在存储引擎内部的相对路径
	SavedPath     string    `gorm:"size:512;not null" json:"-"` 

	// UploadedAt 上传时间
	UploadedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"uploadedAt"`
}

// ListResponse 分页响应包装
type ListResponse struct {
	Items    []AssetRecord `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}