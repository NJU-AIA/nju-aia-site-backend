package livecode

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrDuplicateSlug = errors.New("slug 已存在")
var ErrInvalidBlocks = errors.New("blocks 数据不合法")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateDocument(req UpsertRequest) (string, error) {
	if err := validateBlocks(req.Blocks); err != nil {
		return "", err
	}

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
	}

	blocks := make([]Block, 0, len(req.Blocks))
	for _, block := range req.Blocks {
		blocks = append(blocks, Block{
			ID:         blockID(),
			DocumentID: docID,
			Type:       block.Type,
			Order:      block.Order,
			Content:    block.Content,
			Language:   normalizeLanguage(block.Type, block.Language),
		})
	}

	return docID, s.repo.CreateDocumentWithBlocks(doc, blocks)
}

func (s *Service) GetDocument(id string) (*Document, error) {
	return s.repo.FindDocumentByID(id)
}

func (s *Service) ListDocuments() ([]ListItem, error) {
	return s.repo.ListDocuments()
}

func (s *Service) UpdateDocument(id string, req UpsertRequest) error {
	if err := validateBlocks(req.Blocks); err != nil {
		return err
	}

	doc, err := s.repo.FindDocumentByID(id)
	if err != nil {
		return err
	}

	exists, err := s.repo.ExistsSlug(req.Slug, id)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateSlug
	}

	doc.Slug = req.Slug
	doc.PublishedAt = req.PublishedAt

	blocks := make([]Block, 0, len(req.Blocks))
	for _, block := range req.Blocks {
		blocks = append(blocks, Block{
			ID:         blockID(),
			DocumentID: id,
			Type:       block.Type,
			Order:      block.Order,
			Content:    block.Content,
			Language:   normalizeLanguage(block.Type, block.Language),
		})
	}

	return s.repo.UpdateDocumentWithBlocks(doc, blocks)
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

func validateBlocks(blocks []BlockRequest) error {
	if len(blocks) == 0 {
		return ErrInvalidBlocks
	}

	seenOrders := make(map[int]struct{}, len(blocks))
	for _, block := range blocks {
		if block.Type == "code" && strings.TrimSpace(block.Language) == "" {
			return ErrInvalidBlocks
		}
		if block.Type == "markdown" && strings.TrimSpace(block.Language) != "" {
			return ErrInvalidBlocks
		}
		if _, exists := seenOrders[block.Order]; exists {
			return ErrInvalidBlocks
		}
		seenOrders[block.Order] = struct{}{}
	}

	return nil
}

func normalizeLanguage(blockType string, language string) string {
	if blockType != "code" {
		return ""
	}
	return strings.TrimSpace(language)
}
