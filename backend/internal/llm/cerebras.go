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

// CerebrasClient implements LLMClient using the Cerebras API (OpenAI-compatible).
type CerebrasClient struct {
	APIKey  string
	BaseURL string
	Model   string
	client  *http.Client
}

// NewCerebrasClient creates a new Cerebras client.
func NewCerebrasClient(apiKey, baseURL, model string) *CerebrasClient {
	return &CerebrasClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

// Complete sends a chat completion request to Cerebras.
func (c *CerebrasClient) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	body := chatRequest{
		Model: c.Model,
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
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat/completions", bytes.NewReader(b))
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return CompletionResponse{}, fmt.Errorf("cerebras API error %d: %s", resp.StatusCode, string(respBody))
	}
	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return CompletionResponse{}, fmt.Errorf("parse response: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return CompletionResponse{}, fmt.Errorf("no choices in cerebras response")
	}
	return CompletionResponse{
		Content:    chatResp.Choices[0].Message.Content,
		Model:      chatResp.Model,
		TokensUsed: chatResp.Usage.TotalTokens,
	}, nil
}
