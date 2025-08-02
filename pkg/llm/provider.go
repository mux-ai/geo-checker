package llm

import (
	"context"
	"fmt"
)

type Provider interface {
	Analyze(ctx context.Context, content string, prompt string) (*Response, error)
	Name() string
}

type Response struct {
	Content     string            `json:"content"`
	TokensUsed  int              `json:"tokens_used"`
	Model       string           `json:"model"`
	Metadata    map[string]any   `json:"metadata,omitempty"`
}

type ProviderConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
	BaseURL     string
}

func NewProvider(providerType string, config *ProviderConfig) (Provider, error) {
	switch providerType {
	case "claude":
		return NewClaudeProvider(config)
	case "gpt", "openai":
		return NewOpenAIProvider(config)
	case "local":
		return NewLocalProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}
}