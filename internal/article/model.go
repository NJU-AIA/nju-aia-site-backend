package article

import "time"

type Article struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Category    string    `gorm:"size:100;not null" json:"category"` // 新增分类
	Author      string    `gorm:"size:100;not null" json:"author"`
	DefaultMode string    `gorm:"size:20;not null" json:"defaultMode"` // 枚举存储
	Content     string    `gorm:"type:longtext;not null" json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}