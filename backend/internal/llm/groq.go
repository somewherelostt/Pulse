package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	MaxTokens      int           `json:"max_tokens,omitempty"`
	Temperature    float64       `json:"temperature,omitempty"`
	ResponseFormat interface{}   `json:"response_format,omitempty"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatUsage struct {
	TotalTokens int `json:"total_tokens"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
	Model   string       `json:"model"`
	Usage   chatUsage    `json:"usage"`
}

// GroqClient implements LLMClient using the Groq API (OpenAI-compatible).
type GroqClient struct {
	APIKey  string
	BaseURL string
	Model   string
	client  *http.Client
}

// NewGroqClient creates a new Groq client.
func NewGroqClient(apiKey, baseURL, model string) *GroqClient {
	return &GroqClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

// Complete sends a chat completion request to Groq.
func (g *GroqClient) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	body := chatRequest{
		Model: g.Model,
		Messages: []chatMessage{
			{Role: "system", Content: req.SystemPrompt},
			{Role: "user", Content: req.UserPrompt},
		},
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}
	if req.JSONMode {
		body.ResponseFormat = map[string]string{"type": "json_object"}
	}
	b, err := json.Marshal(body)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("marshal request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.BaseURL+"/chat/completions", bytes.NewReader(b))
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.APIKey)

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return CompletionResponse{}, fmt.Errorf("groq API error %d: %s", resp.StatusCode, string(respBody))
	}
	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return CompletionResponse{}, fmt.Errorf("parse response: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return CompletionResponse{}, fmt.Errorf("no choices in groq response")
	}
	return CompletionResponse{
		Content:    chatResp.Choices[0].Message.Content,
		Model:      chatResp.Model,
		TokensUsed: chatResp.Usage.TotalTokens,
	}, nil
}
