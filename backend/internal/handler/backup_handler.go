package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// BackupHandler 备份处理器
type BackupHandler struct {
	backupService service.BackupService
}

// NewBackupHandler 创建备份处理器
func NewBackupHandler(backupService service.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// GetBackups 获取备份列表
// @Summary 获取备份列表
// @Description 获取所有备份文件列表
// @Tags 数据备份
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/backups [get]
func (h *BackupHandler) GetBackups(c *gin.Context) {
	backups, err := h.backupService.ListBackups()
	if err != nil {
		response.InternalServerError(c, "获取备份列表失败: "+err.Error())
		return
	}

	response.Success(c, backups)
}

// CreateBackup 创建备份
// @Summary 创建备份
// @Description 手动创建数据库备份
// @Tags 数据备份
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "备份创建成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/backups [post]
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	backupInfo, err := h.backupService.CreateBackup()
	if err != nil {
		response.InternalServerError(c, "创建备份失败: "+err.Error())
		return
	}

	response.Success(c, backupInfo)
}

// DownloadBackup 下载备份
// @Summary 下载备份
// @Description 下载指定的备份文件
// @Tags 数据备份
// @Accept json
// @Produce application/gzip
// @Security BearerAuth
// @Param filename path string true "备份文件名"
// @Success 200 {file} file "备份文件"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "备份文件不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/backups/download/{filename} [get]
func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	backupPath, err := h.backupService.GetBackupPath(filename)
	if err != nil {
		response.NotFound(c, "备份文件不存在")
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/gzip")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(backupPath)
}

// DeleteBackup 删除备份
// @Summary 删除备份
// @Description 删除指定的备份文件
// @Tags 数据备份
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param filename path string true "备份文件名"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "备份文件不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/backups/{filename} [delete]
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	err := h.backupService.DeleteBackup(filename)
	if err != nil {
		if err.Error() == "backup file not found" || err.Error() == "invalid filename" {
			response.NotFound(c, "备份文件不存在")
			return
		}
		response.InternalServerError(c, "删除备份失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// CleanupBackups 清理旧备份
// @Summary 清理旧备份
// @Description 清理超出保留数量的旧备份文件
// @Tags 数据备份
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param retention_count query int false "保留数量（默认10）"
// @Success 200 {object} response.Response "清理成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/backups/cleanup [post]
func (h *BackupHandler) CleanupBackups(c *gin.Context) {
	retentionCount := 10 // 默认保留10个备份
	if countStr := c.Query("retention_count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
			retentionCount = count
		}
	}

	err := h.backupService.CleanupOldBackups(retentionCount)
	if err != nil {
		response.InternalServerError(c, "清理备份失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"retention_count": retentionCount,
		"message":         "清理完成",
	})
}
