package livecode

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrDuplicateSlug = errors.New("slug 已存在")
var ErrInvalidBlock = errors.New("block 数据不合法")
var ErrInvalidBlockIDs = errors.New("blockIds 数据不合法")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateDocument(req UpsertRequest) (string, error) {
	exists, err := s.repo.ExistsSlug(req.Slug, "")
	if err != nil {
		return "", err
	}
	if exists {
		return "", ErrDuplicateSlug
	}

	docID := shortID()
	doc := &Document{
		ID:          docID,
		Slug:        req.Slug,
		PublishedAt: req.PublishedAt,
		BlockIDs:    []string{},
	}

	if err := s.repo.CreateDocument(doc); err != nil {
		return "", err
	}

	return docID, nil
}

func (s *Service) GetDocument(id string) (*Document, error) {
	return s.repo.FindDocumentByID(id)
}

func (s *Service) ListDocuments() ([]ListItem, error) {
	return s.repo.ListDocuments()
}

func (s *Service) AddBlock(documentID string, req BlockRequest) (*Block, error) {
	if err := validateBlockRequest(req); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindDocumentByID(documentID); err != nil {
		return nil, err
	}

	return s.repo.AddBlock(documentID, req)
}

func (s *Service) UpdateBlock(documentID, blockID string, req BlockRequest) (*Block, error) {
	if err := validateBlockRequest(req); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindDocumentByID(documentID); err != nil {
		return nil, err
	}

	return s.repo.UpdateBlock(documentID, blockID, req)
}

func (s *Service) DeleteBlock(documentID, blockID string) error {
	if _, err := s.repo.FindDocumentByID(documentID); err != nil {
		return err
	}

	return s.repo.DeleteBlock(documentID, blockID)
}

func (s *Service) UpdateBlockIDs(documentID string, req BlockIDsRequest) error {
	if err := validateBlockIDs(req.BlockIDs); err != nil {
		return err
	}

	if _, err := s.repo.FindDocumentByID(documentID); err != nil {
		return err
	}

	return s.repo.UpdateBlockIDs(documentID, req.BlockIDs)
}

func (s *Service) DeleteDocument(id string) error {
	_, err := s.repo.FindDocumentByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return s.repo.DeleteDocumentByID(id)
}

func shortID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
}

func blockID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
}

func validateBlockRequest(block BlockRequest) error {
	if block.Type == "code" && strings.TrimSpace(block.Language) == "" {
		return ErrInvalidBlock
	}
	if block.Type == "markdown" && strings.TrimSpace(block.Language) != "" {
		return ErrInvalidBlock
	}
	return nil
}

func validateBlockIDs(blockIDs []string) error {
	seen := make(map[string]struct{}, len(blockIDs))
	for _, id := range blockIDs {
		if strings.TrimSpace(id) == "" {
			return ErrInvalidBlockIDs
		}
		if _, ok := seen[id]; ok {
			return ErrInvalidBlockIDs
		}
		seen[id] = struct{}{}
	}
	return nil
}

func normalizeLanguage(blockType string, language string) string {
	if blockType != "code" {
		return ""
	}
	return strings.TrimSpace(language)
}
