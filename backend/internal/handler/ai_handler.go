package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// AIHandler AI handler
type AIHandler struct {
	aiService service.AIService
}

// NewAIHandler creates a new AI handler
func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// TranslateArticle translates an article
func (h *AIHandler) TranslateArticle(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		Summary    string `json:"summary"`
		ProviderID *uint  `json:"provider_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	result, err := h.aiService.Translate(req.Title, req.Content, req.Summary, req.ProviderID)
	if err != nil {
		response.InternalServerError(c, "translation failed: "+err.Error())
		return
	}

	response.Success(c, result)
}

// Chat handles SSE streaming chat
func (h *AIHandler) Chat(c *gin.Context) {
	var req struct {
		ProviderID uint                  `json:"provider_id" binding:"required"`
		Messages   []service.ChatMessage `json:"messages" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	c.Stream(func(w io.Writer) bool {
		err := h.aiService.ChatStream(req.ProviderID, req.Messages, w)
		if err != nil {
			fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
		}
		return false
	})
}

// ListProviders returns all AI providers
func (h *AIHandler) ListProviders(c *gin.Context) {
	providers, err := h.aiService.ListProviders()
	if err != nil {
		response.InternalServerError(c, "failed to list providers")
		return
	}
	response.Success(c, providers)
}

// GetProvider returns a single AI provider
func (h *AIHandler) GetProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	provider, err := h.aiService.GetProvider(uint(id))
	if err != nil {
		response.NotFound(c, "provider not found")
		return
	}
	response.Success(c, provider)
}

// CreateProvider creates a new AI provider
func (h *AIHandler) CreateProvider(c *gin.Context) {
	var req service.CreateAIProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}
	provider, err := h.aiService.CreateProvider(&req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, "created", provider)
}

// UpdateProvider updates an AI provider
func (h *AIHandler) UpdateProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	var req service.UpdateAIProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}
	provider, err := h.aiService.UpdateProvider(uint(id), &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, provider)
}

// DeleteProvider deletes an AI provider
func (h *AIHandler) DeleteProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if err := h.aiService.DeleteProvider(uint(id)); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.NoContent(c, "deleted")
}

// CheckProvider checks provider connectivity
func (h *AIHandler) CheckProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	result, err := h.aiService.CheckProvider(uint(id))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}
