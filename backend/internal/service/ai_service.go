package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/repository"
)

var (
	ErrNoEnabledProvider = errors.New("no enabled AI provider")
	ErrTranslationFailed = errors.New("translation failed")
)

// AIService AI service interface
type AIService interface {
	Translate(title, content, summary string, providerID *uint) (*TranslateResult, error)
	ChatStream(providerID uint, messages []ChatMessage, writer io.Writer) error
	CheckProvider(id uint) (*ProviderCheckResult, error)
	ListProviders() ([]*models.AIProvider, error)
	GetProvider(id uint) (*models.AIProvider, error)
	CreateProvider(req *CreateAIProviderRequest) (*models.AIProvider, error)
	UpdateProvider(id uint, req *UpdateAIProviderRequest) (*models.AIProvider, error)
	DeleteProvider(id uint) error
}

// TranslateResult translation result
type TranslateResult struct {
	TitleEn   string `json:"title_en"`
	ContentEn string `json:"content_en"`
	SummaryEn string `json:"summary_en"`
}

// ChatMessage chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ProviderCheckResult provider connectivity check result
type ProviderCheckResult struct {
	Available bool    `json:"available"`
	Balance   float64 `json:"balance"`
	Error     string  `json:"error,omitempty"`
}

// CreateAIProviderRequest create provider request
type CreateAIProviderRequest struct {
	Name         string `json:"name" binding:"required"`
	ProviderType string `json:"provider_type" binding:"required"`
	APIKey       string `json:"api_key" binding:"required"`
	BaseURL      string `json:"base_url"`
	Model        string `json:"model" binding:"required"`
	IsEnabled    bool   `json:"is_enabled"`
	SortOrder    int    `json:"sort_order"`
}

// UpdateAIProviderRequest update provider request
type UpdateAIProviderRequest struct {
	Name      *string `json:"name"`
	APIKey    *string `json:"api_key"`
	BaseURL   *string `json:"base_url"`
	Model     *string `json:"model"`
	IsEnabled *bool   `json:"is_enabled"`
	SortOrder *int    `json:"sort_order"`
}

type aiService struct {
	providerRepo repository.AIProviderRepository
	crypt        *crypto.Crypto
}

// NewAIService creates a new AI service
func NewAIService(providerRepo repository.AIProviderRepository, masterKey string) (AIService, error) {
	c, err := crypto.NewCrypto(masterKey)
	if err != nil {
		return nil, err
	}
	return &aiService{providerRepo: providerRepo, crypt: c}, nil
}

// Translate translates article content using enabled AI providers
func (s *aiService) Translate(title, content, summary string, providerID *uint) (*TranslateResult, error) {
	providers, err := s.providerRepo.FindEnabled()
	if err != nil || len(providers) == 0 {
		return nil, ErrNoEnabledProvider
	}

	var lastErr error
	for _, p := range providers {
		if providerID != nil && p.ID != *providerID {
			continue
		}
		result, err := s.translateWithProvider(p, title, content, summary)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, fmt.Errorf("%w: %v", ErrTranslationFailed, lastErr)
	}
	return nil, ErrTranslationFailed
}

func (s *aiService) translateWithProvider(provider *models.AIProvider, title, content, summary string) (*TranslateResult, error) {
	apiKey, err := s.crypt.Decrypt(provider.APIKey)
	if err != nil {
		apiKey = provider.APIKey
	}

	prompt := fmt.Sprintf(`Translate the following blog article from Chinese to English. Return ONLY valid JSON with keys: title_en, content_en, summary_en.

Chinese Title: %s
Chinese Summary: %s
Chinese Content (HTML, preserve HTML tags):
%s`, title, summary, content)

	baseURL := s.getBaseURL(provider)
	reqBody := map[string]interface{}{
		"model": provider.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional translator. Translate Chinese to natural, fluent English. Preserve all HTML tags. Return ONLY valid JSON."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	switch provider.ProviderType {
	case "claude":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, errors.New("empty response from AI")
	}

	content = result.Choices[0].Message.Content
	content = extractJSON(content)

	var translateResult TranslateResult
	if err := json.Unmarshal([]byte(content), &translateResult); err != nil {
		return nil, fmt.Errorf("failed to parse translation result: %w", err)
	}

	return &translateResult, nil
}

// ChatStream streams chat response via SSE
func (s *aiService) ChatStream(providerID uint, messages []ChatMessage, writer io.Writer) error {
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return err
	}

	apiKey, _ := s.crypt.Decrypt(provider.APIKey)
	if apiKey == "" {
		apiKey = provider.APIKey
	}

	baseURL := s.getBaseURL(provider)

	msgs := make([]map[string]string, len(messages))
	for i, m := range messages {
		msgs[i] = map[string]string{"role": m.Role, "content": m.Content}
	}

	reqBody := map[string]interface{}{
		"model":    provider.Model,
		"messages": msgs,
		"stream":   true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	switch provider.ProviderType {
	case "claude":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				fmt.Fprintf(writer, "data: [DONE]\n\n")
				break
			}
			fmt.Fprintf(writer, "data: %s\n\n", data)
		}
	}

	return nil
}

// CheckProvider checks provider connectivity and balance
func (s *aiService) CheckProvider(id uint) (*ProviderCheckResult, error) {
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	apiKey, _ := s.crypt.Decrypt(provider.APIKey)
	if apiKey == "" {
		apiKey = provider.APIKey
	}

	baseURL := s.getBaseURL(provider)
	reqBody := map[string]interface{}{
		"model":      provider.Model,
		"messages":   []map[string]string{{"role": "user", "content": "ping"}},
		"max_tokens": 1,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	switch provider.ProviderType {
	case "claude":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &ProviderCheckResult{Available: false, Error: err.Error()}, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	result := &ProviderCheckResult{Available: resp.StatusCode == 200}
	if resp.StatusCode != 200 {
		result.Error = string(body)
	}

	now := time.Now()
	provider.LastCheckAt = &now
	_ = s.providerRepo.Update(provider)

	return result, nil
}

func (s *aiService) getBaseURL(provider *models.AIProvider) string {
	if provider.BaseURL != "" {
		return provider.BaseURL
	}
	switch provider.ProviderType {
	case "claude":
		return "https://api.anthropic.com/v1/messages"
	case "openai":
		return "https://api.openai.com/v1/chat/completions"
	case "deepseek":
		return "https://api.deepseek.com/v1/chat/completions"
	default:
		return provider.BaseURL
	}
}

// ListProviders returns all AI providers
func (s *aiService) ListProviders() ([]*models.AIProvider, error) {
	return s.providerRepo.FindAll()
}

// GetProvider returns a single AI provider
func (s *aiService) GetProvider(id uint) (*models.AIProvider, error) {
	return s.providerRepo.FindByID(id)
}

// CreateProvider creates a new AI provider
func (s *aiService) CreateProvider(req *CreateAIProviderRequest) (*models.AIProvider, error) {
	encryptedKey, err := s.crypt.Encrypt(req.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt API key: %w", err)
	}

	provider := &models.AIProvider{
		Name:         req.Name,
		ProviderType: req.ProviderType,
		APIKey:       encryptedKey,
		BaseURL:      req.BaseURL,
		Model:        req.Model,
		IsEnabled:    req.IsEnabled,
		SortOrder:    req.SortOrder,
	}

	if err := s.providerRepo.Create(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

// UpdateProvider updates an existing AI provider
func (s *aiService) UpdateProvider(id uint, req *UpdateAIProviderRequest) (*models.AIProvider, error) {
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		provider.Name = *req.Name
	}
	if req.APIKey != nil {
		encrypted, err := s.crypt.Encrypt(*req.APIKey)
		if err != nil {
			return nil, err
		}
		provider.APIKey = encrypted
	}
	if req.BaseURL != nil {
		provider.BaseURL = *req.BaseURL
	}
	if req.Model != nil {
		provider.Model = *req.Model
	}
	if req.IsEnabled != nil {
		provider.IsEnabled = *req.IsEnabled
	}
	if req.SortOrder != nil {
		provider.SortOrder = *req.SortOrder
	}

	if err := s.providerRepo.Update(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

// DeleteProvider deletes an AI provider (soft delete)
func (s *aiService) DeleteProvider(id uint) error {
	return s.providerRepo.Delete(id)
}

// extractJSON extracts JSON from AI response text (handles markdown code blocks)
func extractJSON(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}
	return content
}
