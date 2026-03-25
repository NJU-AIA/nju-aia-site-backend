package blog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// CreateRequest 发布文章的请求体
type CreateRequest struct {
	Title   string `json:"title" binding:"required" example:"关于后端严谨性的深度思考"`
	Content string `json:"content" binding:"required" example:"# 第一章：拒绝模糊逻辑\n这是一篇关于 Go 强类型的文章。"`
	Author  string `json:"author" example:"社团负责人"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Message string `json:"error" example:"参数校验失败：标题不能为空"`
}

// CreateBlog 发布文章
// @Summary 发布一篇 Markdown 博客
// @Tags 博客接口
// @Accept json
// @Produce json
// @Param blog body CreateRequest true "文章数据"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /blogs [post]
func (h *Handler) CreateBlog(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "标题和内容为必填项"})
		return
	}
	id, err := h.svc.SaveMarkdown(req.Title, req.Author, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "数据库写入失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "发布成功"})
}

// GetBlog 获取文章
// @Summary 获取特定 ID 的博客
// @Tags 博客接口
// @Param id path int true "博客ID" example(1) 
// @Success 200 {object} Blog
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /blogs/{id} [get]
func (h *Handler) GetBlog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "ID 必须是正整数"})
		return
	}

	b, err := h.svc.GetBlog(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, b)
}

// ListBlogs 获取文章列表
// @Summary 获取所有博客文章列表
// @Tags 博客接口
// @Produce json
// @Success 200 {array} Blog
// @Router /blogs [get]
func (h *Handler) ListBlogs(c *gin.Context) {
	blogs, err := h.svc.repo.List() // 直接调用 repo 获取列表
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	c.JSON(http.StatusOK, blogs)
}