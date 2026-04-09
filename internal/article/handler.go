package article

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 负责文章模块的 HTTP 路由处理
type Handler struct {
	svc *Service
}

// NewHandler 初始化处理器
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ErrorResponse 统一错误响应结构
type ErrorResponse struct {
	Error string `json:"error" example:"错误信息描述"`
}

// ArticleListResponse 匹配前端 Response 接口
type ArticleListResponse struct {
	Items []Article `json:"items"`
}

// CreateArticle 处理发布文章请求
// @Summary 发布文章
// @Tags Articles
// @Accept json
// @Produce json
// @Param body body CreateRequest true "文章数据"
// @Success 201 {object} map[string]string "{"id": "uuid"}"
// @Failure 400 {object} ErrorResponse
// @Router /articles [post]
func (h *Handler) CreateArticle(c *gin.Context) {
	var req CreateRequest
	// 1. 绑定并校验 JSON 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败，请检查必填项、模式和日期格式是否正确"})
		return
	}

	// 2. 调用业务逻辑
	id, err := h.svc.CreateArticle(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "服务器内部错误，保存失败"})
		return
	}

	// 3. 返回 201 Created
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ListArticles 处理获取文章列表请求
// @Summary 获取文章列表
// @Tags Articles
// @Produce json
// @Success 200 {object} ArticleListResponse
// @Router /articles [get]
func (h *Handler) ListArticles(c *gin.Context) {
	list, err := h.svc.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法获取文章列表"})
		return
	}

	c.JSON(http.StatusOK, ArticleListResponse{Items: list})
}

// GetArticle 处理获取单篇文章详情请求
// @Summary 获取文章详情
// @Tags Articles
// @Param id path string true "文章 ID (UUID)"
// @Produce json
// @Success 200 {object} Article
// @Failure 404 {object} ErrorResponse
// @Router /articles/{id} [get]
func (h *Handler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	art, err := h.svc.GetArticle(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "该文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "查询文章失败"})
		}
		return
	}

	c.JSON(http.StatusOK, art)
}

// UpdateArticle 处理编辑文章请求
// @Summary 编辑文章
// @Tags Articles
// @Param id path string true "文章 ID"
// @Param body body CreateRequest true "更新的数据"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /articles/{id} [put]
func (h *Handler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数格式错误，请检查日期格式是否为 YYYY-MM-DD"})
		return
	}

	if err := h.svc.UpdateArticle(id, req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "无法更新，文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "更新操作失败"})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteArticle 处理删除文章请求
// @Summary 删除文章
// @Tags Articles
// @Param id path string true "文章 ID"
// @Success 204 "No Content"
// @Router /articles/{id} [delete]
func (h *Handler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "删除操作失败"})
		return
	}

	c.Status(http.StatusNoContent)
}
