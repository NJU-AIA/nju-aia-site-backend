package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// CreateRequest 对应 Apifox Body 参数
type CreateRequest struct {
	Title    string `json:"title" binding:"required" example:"关于后端严谨性的深度思考"`
	Category string `json:"category" binding:"required" example:"技术分享"`
	Author   string `json:"author" binding:"required" example:"负责人"`
	// @Enums article, slide, homework
	DefaultMode string `json:"defaultMode" binding:"required,oneof=article slide homework" example:"article"`
	Content     string `json:"content" binding:"required" example:"# Markdown 内容"`
}

// ErrorResponse 统一错误响应
type ErrorResponse struct {
	Error string `json:"error" example:"参数校验失败"`
}

// CreateArticle 发布新文章
// @Summary 发布新文章
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body CreateRequest true "文章内容"
// @Success 201 {object} map[string]string "{"id": "uuid"}"
// @Failure 400 {object} ErrorResponse
// @Router /articles [post]
func (h *Handler) CreateArticle(c *gin.Context) {
	var req CreateRequest
	
	// 校验参数（包含必填项和枚举值校验）
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数校验失败: title, category, author, content, defaultMode 为必填且 defaultMode 必须为规定枚举值"})
		return
	}

	id, err := h.svc.CreateArticle(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "数据库保存失败"})
		return
	}

	// 对应 Apifox 要求：HTTP 201，仅返回 ID
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// GetArticle 获取文章
// @Summary 获取特定 ID 的博客
// @Tags 博客接口
// @Param id path string true "博客ID" example("1") 
// @Success 200 {object} Article
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /articles/{id} [get]
func (h *Handler) GetArticle(c *gin.Context) {
	// 1. 直接获取路径参数（它本身就是 string）
	id := c.Param("id")

	// 2. 直接传给 service，不需要转成 uint
	b, err := h.svc.GetArticle(id) 
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, b)
}

// ListArticles 获取文章列表
// @Summary 获取所有博客文章列表
// @Tags 博客接口
// @Produce json
// @Success 200 {array} Article
// @Router /articles [get]
func (h *Handler) ListArticles(c *gin.Context) {
	articles, err := h.svc.repo.List() // 直接调用 repo 获取列表
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "获取列表失败"})
		return
	}
	c.JSON(http.StatusOK, articles)
}