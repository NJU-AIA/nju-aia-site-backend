package asset

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

// ErrorResponse 统一错误响应
type ErrorResponse struct {
	Error string `json:"error" example:"错误描述"`
}

// ListQuery 定义获取列表时的查询参数
type ListQuery struct {
	Scope    string `form:"scope"`
	Kind     string `form:"kind" binding:"omitempty,oneof=article shared"`
	Keyword  string `form:"keyword"`
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize,default=20" binding:"omitempty,min=1,max=100"`
}

// UploadFile 处理静态资源上传
// @Summary 上传静态资源
// @Description 按照 scope 分类存储，支持文章私有/{articleId}/filename 和公共/{scope}/filename
// @Tags Assets
// @Accept multipart/form-data
// @Produce json
// @Param scope formData string true "第一层路径：文章ID或固定目录(images, videos等)"
// @Param name formData string true "文件名主干（不含扩展名）"
// @Param overwrite formData bool false "是否允许覆盖同路径资源"
// @Param file formData file true "要上传的文件"
// @Success 201 {object} AssetRecord
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "资源已存在"
// @Router /assets [post]
func (h *Handler) UploadFile(c *gin.Context) {
	// 1. 获取表单参数
	scope := c.PostForm("scope")
	nameStem := c.PostForm("name")
	overwrite := c.PostForm("overwrite") == "true" // 字符串转布尔

	// 2. 获取文件对象
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "未接收到有效文件"})
		return
	}

	// 基础必填校验
	if scope == "" || nameStem == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "scope 和 name 为必填项"})
		return
	}

	// 3. 调用 Service 进行逻辑预处理
	// record: 数据库模型数据, savedPath: 存储引擎所需的物理路径(不带前缀斜杠)
	record, savedPath, err := h.svc.ProcessUpload(fileHeader.Filename, scope, nameStem)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	record.Size = fileHeader.Size

	// 4. 检查路径冲突 (唯一索引保护)
	// 根据生成的逻辑路径（如 /images/capoo.jpg）查询数据库
	existing, err := h.svc.repo.FindByPath(record.Path)
	if err == nil && existing != nil {
		if !overwrite {
			// 如果已存在且未开启覆盖模式，返回 409 Conflict
			c.JSON(http.StatusConflict, ErrorResponse{Error: "资源已存在，如需覆盖请开启 overwrite 模式"})
			return
		}
		// 开启了覆盖：将旧记录的 ID 赋给新记录，后续执行 Update 而非 Create
		record.ID = existing.ID
	}

	// 5. 执行物理保存 (通过注入的 Storage 接口，支持本地/COS)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法读取文件流"})
		return
	}
	defer file.Close()

	if err := h.svc.storage.Save(savedPath, file); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "物理存储写入失败: " + err.Error()})
		return
	}

	// 6. 数据库持久化
	if record.ID > 0 {
		// 更新已有记录
		if err := h.svc.repo.Update(record); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "更新数据库记录失败"})
			return
		}
	} else {
		// 创建新记录
		if err := h.svc.repo.Create(record); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "创建数据库记录失败"})
			return
		}
	}

	// 7. 返回 201 Created 和生成的 AssetRecord (包含小驼峰字段及正确的 markdownValue)
	c.JSON(http.StatusCreated, record)
}

// ListAssets 获取资源列表
// @Summary 获取静态资源列表
// @Tags Assets
// @Produce json
// @Param scope query string false "路径范围"
// @Param kind query string false "类型" enums(article,shared)
// @Param page query int false "页码"
// @Router /assets [get]
func (h *Handler) ListAssets(c *gin.Context) {
	var q ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "查询参数格式错误"})
		return
	}

	items, total, err := h.svc.repo.List(q.Scope, q.Kind, q.Keyword, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Items:    items,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	})
}

// DeleteAsset 删除静态资源
// @Summary 删除静态资源
// @Tags Assets
// @Param path query string true "资源完整路径 (如 /images/capoo.webp)"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /assets [delete]
func (h *Handler) DeleteAsset(c *gin.Context) {
	// 获取参数: /storage/assets/images/capoo.jpg
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "必须提供 path 参数"})
		return
	}

	// 1. 数据库查询 (必须精确匹配 /storage/assets/...)
	record, err := h.svc.repo.FindByPath(path)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// 2. 物理删除
	_ = h.svc.storage.Delete(record.SavedPath)

	// 3. 数据库删除
	if err := h.svc.repo.Delete(record.ID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "记录删除失败"})
		return
	}

	c.Status(http.StatusNoContent)
}