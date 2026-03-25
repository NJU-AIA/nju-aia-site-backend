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

type SuccessResponse struct {
	URL string `json:"url"`
	ID  uint   `json:"id"`
}

// @Summary 上传资源
// @Tags 资源接口
// @Accept multipart/form-data
// @Param file formData file true "文件"
// @Success 200 {object} SuccessResponse
// @Router /assets [post]
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件接收失败"})
		return
	}

	pPath, aURL, fType, err := h.svc.GeneratePath(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "路径生成失败"})
		return
	}

	if err := c.SaveUploadedFile(file, pPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}

	record := &AssetRecord{
		FileName: file.Filename, SavedPath: pPath,
		AccessURL: aURL, FileType: fType,
	}
	h.svc.repo.Create(record)

	c.JSON(http.StatusOK, SuccessResponse{URL: aURL, ID: record.ID})
}