package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iambaby/blog/internal/pkg/response"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	uploadDir string // 上传目录
	maxSize   int64  // 最大文件大小（字节）
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(uploadDir string, maxSizeMB int64) *UploadHandler {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}

	return &UploadHandler{
		uploadDir: uploadDir,
		maxSize:   maxSizeMB * 1024 * 1024, // 转换为字节
	}
}

// UploadImageResponse 上传图片响应
type UploadImageResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// UploadImage 上传图片
// @Summary 上传图片
// @Description 上传文章图片
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "图片文件"
// @Success 200 {object} response.Response{data=UploadImageResponse} "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 413 {object} response.Response "文件太大"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 检查文件大小
	if file.Size > h.maxSize {
		response.BadRequest(c, fmt.Sprintf("文件大小不能超过 %dMB", h.maxSize/1024/1024))
		return
	}

	// 检查文件类型
	if !h.isImage(file) {
		response.BadRequest(c, "只支持图片文件（jpg, jpeg, png, gif, webp）")
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 创建日期子目录（按年月组织）
	dateDir := time.Now().Format("2006/01")
	targetDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		response.InternalServerError(c, "创建目录失败")
		return
	}

	// 保存文件
	targetPath := filepath.Join(targetDir, filename)
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		response.InternalServerError(c, "保存文件失败: "+err.Error())
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/uploads/%s/%s", dateDir, filename)

	response.SuccessWithMessage(c, "上传成功", UploadImageResponse{
		URL:      fileURL,
		Filename: file.Filename,
		Size:     file.Size,
	})
}

// UploadArticleImage 上传文章编辑器中的图片
// @Summary 上传文章编辑器图片
// @Description 上传文章编辑器中的图片（专用于富文本编辑器）
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "图片文件"
// @Success 200 {object} response.Response{data=UploadImageResponse} "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 413 {object} response.Response "文件太大"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/upload/article-image [post]
func (h *UploadHandler) UploadArticleImage(c *gin.Context) {
	// 和 UploadImage 相同，但可以增加额外的逻辑
	h.UploadImage(c)
}

// isImage 检查是否为图片文件
func (h *UploadHandler) isImage(file *multipart.FileHeader) bool {
	// 检查扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return false
	}

	// 检查MIME类型
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	// 读取文件头（前512字节）
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// 检测MIME类型
	contentType := ""
	if len(buffer) > 0 {
		// 简单的MIME类型检测
		if len(buffer) >= 2 {
			// JPEG
			if buffer[0] == 0xFF && buffer[1] == 0xD8 {
				contentType = "image/jpeg"
			}
		}
		if len(buffer) >= 8 {
			// PNG
			if buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47 {
				contentType = "image/png"
			}
		}
		if len(buffer) >= 6 {
			// GIF
			if buffer[0] == 0x47 && buffer[1] == 0x49 && buffer[2] == 0x46 {
				contentType = "image/gif"
			}
		}
		if len(buffer) >= 12 {
			// WebP
			if buffer[0] == 0x52 && buffer[1] == 0x49 && buffer[2] == 0x46 && buffer[3] == 0x46 &&
				buffer[8] == 0x57 && buffer[9] == 0x45 && buffer[10] == 0x42 && buffer[11] == 0x50 {
				contentType = "image/webp"
			}
		}
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	return allowedTypes[contentType]
}
