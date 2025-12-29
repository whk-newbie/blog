package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/jwt"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// ConfigHandler 配置处理器
type ConfigHandler struct {
	configService service.ConfigService
}

// NewConfigHandler 创建配置处理器
func NewConfigHandler(configService service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// GetConfigs 获取配置列表
// @Summary 获取配置列表
// @Description 获取系统配置列表
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param config_type query string false "配置类型"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/configs [get]
func (h *ConfigHandler) GetConfigs(c *gin.Context) {
	configType := c.Query("config_type")

	configs, err := h.configService.GetConfigs(configType)
	if err != nil {
		response.InternalServerError(c, "获取配置列表失败: "+err.Error())
		return
	}

	response.Success(c, configs)
}

// GetConfigByID 获取配置详情
// @Summary 获取配置详情
// @Description 根据ID获取配置详情
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "配置不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/configs/:id [get]
func (h *ConfigHandler) GetConfigByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的配置ID")
		return
	}

	config, err := h.configService.GetConfigByID(uint(id))
	if err != nil {
		if err == service.ErrConfigNotFound {
			response.NotFound(c, "配置不存在")
			return
		}
		response.InternalServerError(c, "获取配置失败: "+err.Error())
		return
	}

	response.Success(c, config)
}

// CreateConfig 创建配置
// @Summary 创建配置
// @Description 创建新的系统配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateConfigRequest true "创建配置请求"
// @Success 201 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 409 {object} response.Response "配置已存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/configs [post]
func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	var req service.CreateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := claims.(*jwt.Claims).UserID

	config, err := h.configService.CreateConfig(&req, userID)
	if err != nil {
		if err == service.ErrConfigExists {
			response.Error(c, 409, "配置已存在")
			return
		}
		response.InternalServerError(c, "创建配置失败: "+err.Error())
		return
	}

	// 创建成功返回201状态码
	c.JSON(201, gin.H{
		"code":    0,
		"message": "配置创建成功",
		"data":    config,
	})
}

// UpdateConfig 更新配置
// @Summary 更新配置
// @Description 更新系统配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Param request body service.UpdateConfigRequest true "更新配置请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "配置不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/configs/:id [put]
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的配置ID")
		return
	}

	var req service.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := claims.(*jwt.Claims).UserID

	config, err := h.configService.UpdateConfig(uint(id), &req, userID)
	if err != nil {
		if err == service.ErrConfigNotFound {
			response.NotFound(c, "配置不存在")
			return
		}
		response.InternalServerError(c, "更新配置失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "配置更新成功", config)
}

// DeleteConfig 删除配置
// @Summary 删除配置
// @Description 删除系统配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "配置不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/configs/:id [delete]
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的配置ID")
		return
	}

	err = h.configService.DeleteConfig(uint(id))
	if err != nil {
		if err == service.ErrConfigNotFound {
			response.NotFound(c, "配置不存在")
			return
		}
		response.InternalServerError(c, "删除配置失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "配置删除成功", nil)
}

// GenerateCrawlerToken 生成爬虫Token
// @Summary 生成爬虫Token
// @Description 生成新的爬虫认证Token
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "生成Token请求" example({"name": "爬虫Token #1"})
// @Success 200 {object} response.Response "生成成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/configs/generate-crawler-token [post]
func (h *ConfigHandler) GenerateCrawlerToken(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := claims.(*jwt.Claims).UserID

	token, err := h.configService.GenerateCrawlerToken(req.Name, userID)
	if err != nil {
		response.InternalServerError(c, "生成Token失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Token生成成功", token)
}
