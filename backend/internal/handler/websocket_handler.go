package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gorillaWS "github.com/gorilla/websocket"
	"github.com/whk-newbie/blog/internal/pkg/jwt"
	"github.com/whk-newbie/blog/internal/websocket"
)

var upgrader = gorillaWS.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源（生产环境应该限制）
		return true
	},
}

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	hub        *websocket.Hub
	jwtManager *jwt.Manager
}

// NewWebSocketHandler 创建WebSocket处理器
func NewWebSocketHandler(hub *websocket.Hub, jwtManager *jwt.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		jwtManager: jwtManager,
	}
}

// HandleCrawlerTasks 处理爬虫任务WebSocket连接
// @Summary WebSocket连接
// @Description 建立WebSocket连接，接收爬虫任务实时更新
// @Tags 爬虫任务
// @Accept json
// @Produce json
// @Param token query string true "JWT Token"
// @Success 101 {object} nil "WebSocket连接成功"
// @Failure 400 {object} github_com_whk-newbie_blog_internal_pkg_response.Response "参数错误"
// @Failure 401 {object} github_com_whk-newbie_blog_internal_pkg_response.Response "未授权"
// @Router /ws/crawler/tasks [get]
func (h *WebSocketHandler) HandleCrawlerTasks(c *gin.Context) {
	// 从查询参数获取Token
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Token"})
		return
	}

	// 验证Token
	claims, err := h.jwtManager.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效"})
		return
	}

	// 升级为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket升级失败"})
		return
	}

	// 创建客户端
	client := websocket.NewClient(h.hub, conn)

	// 注册客户端
	h.hub.Register(client)

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()

	// 记录用户信息（可选）
	_ = claims
}
