package asset

import "time"

type AssetRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileName  string    `gorm:"size:255" json:"file_name"`
	SavedPath string    `gorm:"size:512" json:"-"` // 隐藏物理路径
	AccessURL string    `gorm:"size:512" json:"url"`
	FileType  string    `gorm:"size:50" json:"file_type"`
	CreatedAt time.Time `json:"createdAt"`
}