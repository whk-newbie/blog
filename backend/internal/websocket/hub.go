package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/whk-newbie/blog/internal/models"
)

const (
	// 写超时时间
	writeWait = 10 * time.Second

	// 读超时时间
	pongWait = 60 * time.Second

	// ping周期（必须小于pongWait）
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512
)

// Hub 维护所有活跃的客户端连接和广播消息
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool

	// 广播消息通道
	broadcast chan []byte

	// 注册客户端通道
	register chan *Client

	// 注销客户端通道
	unregister chan *Client

	// 互斥锁
	mu sync.RWMutex
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastTaskUpdate 广播任务更新
func (h *Hub) BroadcastTaskUpdate(task *models.CrawlTask) {
	message := TaskUpdateMessage{
		Type: "task_update",
		Data: TaskUpdateData{
			TaskID:    task.TaskID,
			Status:    string(task.Status),
			Progress:  task.Progress,
			Message:   task.Message,
			UpdatedAt: task.UpdatedAt,
		},
	}

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.broadcast <- data
}

// Client 表示一个WebSocket客户端连接
type Client struct {
	hub *Hub

	// WebSocket连接
	conn *websocket.Conn

	// 发送消息的缓冲通道
	send chan []byte
}

// NewClient 创建新的客户端
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

// ReadPump 从WebSocket连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 日志记录错误
			}
			break
		}

		// 处理ping消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msg["type"] == "ping" {
				// 响应pong
				pongMsg := map[string]string{"type": "pong"}
				pongData, _ := json.Marshal(pongMsg)
				c.send <- pongData
			}
		}
	}
}

// WritePump 向WebSocket连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量发送队列中的消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// TaskUpdateMessage 任务更新消息
type TaskUpdateMessage struct {
	Type string         `json:"type"`
	Data TaskUpdateData `json:"data"`
}

// TaskUpdateData 任务更新数据
type TaskUpdateData struct {
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updated_at"`
}
