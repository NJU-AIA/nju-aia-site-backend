package asset

import (
	"mime"
	"path/filepath"
	"strings"
)

type Service struct {
	repo    *Repository
	storage Storage
}

func NewService(repo *Repository, storage Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

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

	kind := "article"
	commonScopes := map[string]bool{"images": true, "videos": true, "docs": true, "others": true}
	if commonScopes[scope] {
		kind = "shared"
	}

	path := "/" + scope + "/" + filename
	savedPath := scope + "/" + filename
	baseURL := ""
	if s.storage != nil {
		baseURL = strings.TrimSuffix(s.storage.GetBaseURL(), "/")
	}
	fullURL := baseURL + path

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
		SavedPath:     savedPath,
	}

	return record, savedPath, nil
}
