package asset

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GeneratePath(originalName string) (pPath, aURL, fType string, err error) {
	ext := strings.ToLower(filepath.Ext(originalName))
	fType = "others"
	if isImage(ext) { fType = "images" } else if isDoc(ext) { fType = "docs" }

	now := time.Now()
	dateDir := now.Format("2006/01/02")
	
	// 通用规范：根目录下的 storage/assets
	relDir := filepath.Join("storage", "assets", fType, dateDir)
	_ = os.MkdirAll(relDir, 0755)

	// 命名：UUID_时间戳.后缀
	newName := fmt.Sprintf("%s_%d%s", uuid.New().String(), now.UnixNano(), ext)
	pPath = filepath.Join(relDir, newName)
	aURL = fmt.Sprintf("/assets/%s/%s/%s", fType, dateDir, newName)

	return pPath, aURL, fType, nil
}

func isImage(ext string) bool {
	return map[string]bool{".jpg":true, ".jpeg":true, ".png":true, ".gif":true, ".webp":true}[ext]
}

func isDoc(ext string) bool {
	return map[string]bool{".pdf":true, ".docx":true, ".txt":true, ".md":true}[ext]
}