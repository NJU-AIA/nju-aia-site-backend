package upload

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

type UploadResponse struct {
	ID      uint   `json:"id" example:"1"`
	URL     string `json:"url" example:"/uploads/images/2026/03/25/uuid.png"`
	Message string `json:"message" example:"上传成功"`
}

type ErrorResponse struct {
	Message string `json:"error" example:"错误描述"`
}

// @Summary 上传资源文件
// @Tags 资源接口
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} ErrorResponse
// @Router /upload [post]
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效文件"})
		return
	}

	pPath, aURL, fType, err := h.svc.GeneratePath(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "目录创建失败"})
		return
	}

	if err := c.SaveUploadedFile(file, pPath); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "保存失败"})
		return
	}

	record := &FileRecord{
		FileName: file.Filename, SavedPath: pPath,
		AccessURL: aURL, FileType: fType, Size: file.Size,
	}
	h.svc.repo.Create(record)

	c.JSON(http.StatusOK, UploadResponse{
		ID: record.ID, URL: aURL, Message: "上传成功",
	})
}