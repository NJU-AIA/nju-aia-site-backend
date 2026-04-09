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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败，请检查发布日期和 blocks 格式"})
		return
	}

	id, err := h.svc.CreateDocument(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateSlug):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "slug 已存在"})
		case errors.Is(err, ErrInvalidBlocks):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "blocks 数据不合法"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "创建 livecode 失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
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

// UpdateDocument 更新 livecode 文件
// @Summary 更新 livecode 文件
// @Tags Livecodes
// @Accept json
// @Param id path string true "livecode ID"
// @Param body body UpsertRequest true "更新的数据"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /livecodes/{id} [put]
func (h *Handler) UpdateDocument(c *gin.Context) {
	var req UpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败，请检查发布日期和 blocks 格式"})
		return
	}

	err := h.svc.UpdateDocument(c.Param("id"), req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "无法更新，livecode 不存在"})
		case errors.Is(err, ErrDuplicateSlug):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "slug 已存在"})
		case errors.Is(err, ErrInvalidBlocks):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "blocks 数据不合法"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "更新 livecode 失败"})
		}
		return
	}

	c.Status(http.StatusOK)
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
