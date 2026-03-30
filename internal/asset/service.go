package asset

import (
	"mime"
	"path/filepath"
	"strings"
)

type Service struct {
	repo    *Repository
	storage Storage // 核心：注入接口
}

func NewService(repo *Repository, storage Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

// isSafeSegment 检查路径分段是否安全：不含 .. / \ 等路径穿越字符。
func isSafeSegment(s string) bool {
	if s == "" || s == "." || s == ".." {
		return false
	}
	return !strings.ContainsAny(s, `/\`)
}

func (s *Service) ProcessUpload(originalName, scope, nameStem string) (*AssetRecord, string, error) {
	if !isSafeSegment(scope) || !isSafeSegment(nameStem) {
		return nil, "", ErrUnsafePath
	}

	ext := filepath.Ext(originalName) 
	filename := nameStem + ext        

	// 1. 判断归属类型
	kind := "article"
	commonScopes := map[string]bool{"images": true, "videos": true, "docs": true, "others": true}
	if commonScopes[scope] {
		kind = "shared"
	}

	// 2. 生成逻辑 Path (用于 API 返回和数据库检索): /{scope}/{filename}
	// 示例: /images/capoo.jpg 或 /article-123/cover.png
	path := "/" + scope + "/" + filename

	// 3. 生成存储 Key (savedPath): {scope}/{filename} 
	// 注意：这里去掉了所有前缀，直接以 scope 开头，确保存在 COS 根目录
	savedPath := scope + "/" + filename

	// 4. 生成 URL
	// 如果是 COS，GetBaseURL 返回 "https://bucket-125...myqcloud.com"
	// 拼接后为 "https://bucket-125...myqcloud.com/images/capoo.jpg"
	fullURL := strings.TrimSuffix(s.storage.GetBaseURL(), "/") + path

	// 5. 生成 markdownValue
	markdownValue := path 
	if kind == "article" {
		markdownValue = filename 
	}

	record := &AssetRecord{
		Path:          path,
		URL:           fullURL,
		MarkdownValue: markdownValue,
		Scope:         scope,
		Kind:          kind,
		Filename:      filename,
		Name:          nameStem,
		Ext:           strings.TrimPrefix(ext, "."),
		ContentType:   mime.TypeByExtension(ext),
		SavedPath:     savedPath, // 数据库存入: images/capoo.jpg
	}

	return record, savedPath, nil
}