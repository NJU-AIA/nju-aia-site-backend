package upload

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

	datePart := time.Now().Format("2006/01/02")
	relDir := filepath.Join("internal", "uploads", fType, datePart)
	_ = os.MkdirAll(relDir, 0755)

	uniqueName := uuid.New().String() + ext
	pPath = filepath.Join(relDir, uniqueName)
	aURL = fmt.Sprintf("/uploads/%s/%s/%s", fType, datePart, uniqueName)

	return pPath, aURL, fType, nil
}

func isImage(ext string) bool {
	return map[string]bool{".jpg":true, ".jpeg":true, ".png":true, ".gif":true, ".webp":true}[ext]
}
func isDoc(ext string) bool {
	return map[string]bool{".pdf":true, ".docx":true, ".txt":true, ".md":true}[ext]
}