package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"123456"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpass123"`
}

// Login 登录
// @Summary 管理员登录
// @Description 使用用户名和密码登录，获取JWT Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	loginResp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Unauthorized(c, "用户名或密码错误")
			return
		}
		response.InternalServerError(c, "登录失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "登录成功", loginResp)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权或旧密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 从上下文获取用户ID（由认证中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	err := h.authService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			response.Unauthorized(c, "旧密码错误")
		case service.ErrPasswordTooShort:
			response.BadRequest(c, "新密码长度至少为6位")
		case service.ErrSamePassword:
			response.BadRequest(c, "新密码不能与旧密码相同")
		default:
			response.InternalServerError(c, "修改密码失败: "+err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// VerifyToken 验证Token
// @Summary 验证Token
// @Description 验证JWT Token是否有效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]interface{}} "Token有效"
// @Failure 401 {object} response.Response "Token无效或已过期"
// @Router /auth/verify [get]
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	// 如果能走到这里，说明Token已经通过认证中间件验证
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	response.Success(c, gin.H{
		"user_id":  userID,
		"username": username,
		"valid":    true,
	})
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 使用旧Token获取新Token
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]string} "刷新成功"
// @Failure 401 {object} response.Response "Token无效或已过期"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从请求头获取Token
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, "缺少Token")
		return
	}

	// 移除"Bearer "前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	newToken, err := h.authService.RefreshToken(token)
	if err != nil {
		response.Unauthorized(c, "Token刷新失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": newToken,
	})
}

