package livecode

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type ErrorResponse struct {
	Error string `json:"error" example:"错误信息描述"`
}

// CreateDocument 创建 livecode 文件
// @Summary 创建 livecode 文件
// @Tags Livecodes
// @Accept json
// @Produce json
// @Param body body UpsertRequest true "livecode 数据"
// @Success 201 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Router /livecodes [post]
func (h *Handler) CreateDocument(c *gin.Context) {
	var req UpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败，请检查 slug 和发布日期格式"})
		return
	}

	id, err := h.svc.CreateDocument(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateSlug):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "slug 已存在"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "创建 livecode 失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateBlockIDs 更新 block 列表顺序
// @Summary 更新 block 列表顺序
// @Tags Livecodes
// @Accept json
// @Param id path string true "livecode ID"
// @Param body body BlockIDsRequest true "block id 列表"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /livecodes/{id}/block-ids [put]
func (h *Handler) UpdateBlockIDs(c *gin.Context) {
	var req BlockIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败，请检查 blockIds"})
		return
	}

	if err := h.svc.UpdateBlockIDs(c.Param("id"), req); err != nil {
		switch {
		case errors.Is(err, ErrInvalidBlockIDs):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "blockIds 数据不合法"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "目标 livecode 或 block 不存在"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "更新 block 列表失败"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// AddBlock 为 livecode 文件新增一个 block
// @Summary 新增 block
// @Tags Livecodes
// @Accept json
// @Produce json
// @Param id path string true "livecode ID"
// @Param body body BlockRequest true "block 数据"
// @Success 201 {object} BlockResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /livecodes/{id}/blocks [post]
func (h *Handler) AddBlock(c *gin.Context) {
	var req BlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败，请检查 block 类型和语言"})
		return
	}

	block, err := h.svc.AddBlock(c.Param("id"), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidBlock):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "block 数据不合法"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "该 livecode 不存在"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "新增 block 失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, BlockResponse{
		ID:       block.ID,
		Type:     block.Type,
		Content:  block.Content,
		Language: block.Language,
	})
}

// UpdateBlock 更新单个 block
// @Summary 更新 block
// @Tags Livecodes
// @Accept json
// @Produce json
// @Param id path string true "livecode ID"
// @Param blockId path string true "block ID"
// @Param body body BlockRequest true "block 数据"
// @Success 200 {object} BlockResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /livecodes/{id}/blocks/{blockId} [put]
func (h *Handler) UpdateBlock(c *gin.Context) {
	var req BlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败，请检查 block 类型和语言"})
		return
	}

	block, err := h.svc.UpdateBlock(c.Param("id"), c.Param("blockId"), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidBlock):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "block 数据不合法"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "目标 livecode 或 block 不存在"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "更新 block 失败"})
		}
		return
	}

	c.JSON(http.StatusOK, BlockResponse{
		ID:       block.ID,
		Type:     block.Type,
		Content:  block.Content,
		Language: block.Language,
	})
}

// DeleteBlock 删除单个 block
// @Summary 删除 block
// @Tags Livecodes
// @Param id path string true "livecode ID"
// @Param blockId path string true "block ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /livecodes/{id}/blocks/{blockId} [delete]
func (h *Handler) DeleteBlock(c *gin.Context) {
	if err := h.svc.DeleteBlock(c.Param("id"), c.Param("blockId")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "目标 livecode 或 block 不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "删除 block 失败"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// ListDocuments 获取 livecode 列表
// @Summary 获取 livecode 列表
// @Tags Livecodes
// @Produce json
// @Success 200 {object} ListResponse
// @Router /livecodes [get]
func (h *Handler) ListDocuments(c *gin.Context) {
	items, err := h.svc.ListDocuments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法获取 livecode 列表"})
		return
	}

	c.JSON(http.StatusOK, ListResponse{Items: items})
}

// GetDocument 获取 livecode 详情
// @Summary 获取 livecode 详情
// @Tags Livecodes
// @Produce json
// @Param id path string true "livecode ID"
// @Success 200 {object} Document
// @Failure 404 {object} ErrorResponse
// @Router /livecodes/{id} [get]
func (h *Handler) GetDocument(c *gin.Context) {
	doc, err := h.svc.GetDocument(c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "该 livecode 不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "查询 livecode 失败"})
		}
		return
	}

	c.JSON(http.StatusOK, doc)
}

// DeleteDocument 删除 livecode 文件
// @Summary 删除 livecode 文件
// @Tags Livecodes
// @Param id path string true "livecode ID"
// @Success 204
// @Router /livecodes/{id} [delete]
func (h *Handler) DeleteDocument(c *gin.Context) {
	if err := h.svc.DeleteDocument(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "删除 livecode 失败"})
		return
	}

	c.Status(http.StatusNoContent)
}
