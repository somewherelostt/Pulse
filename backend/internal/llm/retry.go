package llm

import (
	"context"
	"fmt"
	"log/slog"
)

// FallbackClient tries primary first, then falls back to secondary on error.
type FallbackClient struct {
	primary  LLMClient
	fallback LLMClient
}

// NewFallbackClient creates a fallback-aware LLM client.
func NewFallbackClient(primary, fallback LLMClient) *FallbackClient {
	return &FallbackClient{primary: primary, fallback: fallback}
}

// Complete tries the primary provider; on failure, tries the fallback.
func (f *FallbackClient) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	resp, err := f.primary.Complete(ctx, req)
	if err != nil {
		slog.Warn("primary LLM provider failed, switching to fallback", "err", err)
		resp, err = f.fallback.Complete(ctx, req)
		if err != nil {
			return CompletionResponse{}, fmt.Errorf("both LLM providers failed: %w", err)
		}
		resp.Model = "cerebras/" + resp.Model
		slog.Info("fallback LLM provider succeeded", "model", resp.Model)
	}
	return resp, nil
}
