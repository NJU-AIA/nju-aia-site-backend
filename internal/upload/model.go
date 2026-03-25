package upload

import "time"

type FileRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	FileName  string    `gorm:"size:255;not null" json:"file_name" example:"image.png"`
	SavedPath string    `gorm:"size:512;not null" json:"-"` // 物理路径不展示
	AccessURL string    `gorm:"size:512;not null" json:"url" example:"/uploads/images/2026/03/25/uuid.png"`
	FileType  string    `gorm:"size:50" json:"file_type" example:"images"`
	Size      int64     `json:"size" example:"102400"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}