package llm

import "context"

// LLMClient is the common interface for all LLM providers.
type LLMClient interface {
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

// CompletionRequest holds the parameters for a chat completion call.
type CompletionRequest struct {
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
	Temperature  float64
	JSONMode     bool // if true, sets response_format: json_object
}

// CompletionResponse holds the result from a chat completion call.
type CompletionResponse struct {
	Content    string
	Model      string
	TokensUsed int
}
