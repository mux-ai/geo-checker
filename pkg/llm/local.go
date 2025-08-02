package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type LocalProvider struct {
	config *ProviderConfig
	client *http.Client
}

type localRequest struct {
	Model       string        `json:"model"`
	Messages    []message     `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream"`
}

type localResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Model string `json:"model"`
}

func NewLocalProvider(config *ProviderConfig) (*LocalProvider, error) {
	if config == nil {
		return nil, NewLLMError(ErrorTypeRequest, "Provider configuration is required", "local")
	}
	
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}
	
	// Validate URL format
	if _, err := url.Parse(config.BaseURL); err != nil {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("Invalid base URL: %v", err), "local")
	}
	
	if config.Model == "" {
		config.Model = "llama2"
	}
	
	// For local models, we're more permissive with model names
	if strings.TrimSpace(config.Model) == "" {
		return nil, NewLLMError(ErrorTypeModel, "Model name cannot be empty", "local")
	}
	
	// Validate temperature if set
	if config.Temperature < 0 || config.Temperature > 2 {
		return nil, NewLLMError(ErrorTypeRequest, "Temperature must be between 0 and 2", "local")
	}
	
	return &LocalProvider{
		config: config,
		client: &http.Client{Timeout: 120 * time.Second},
	}, nil
}

func (l *LocalProvider) Name() string {
	return "local"
}

func (l *LocalProvider) Analyze(ctx context.Context, content string, prompt string) (*Response, error) {
	// Validate inputs
	if strings.TrimSpace(content) == "" {
		return nil, NewLLMError(ErrorTypeRequest, "Content cannot be empty - webpage scraping may have failed or returned no extractable content", "local")
	}
	if strings.TrimSpace(prompt) == "" {
		return nil, NewLLMError(ErrorTypeRequest, "Prompt cannot be empty", "local")
	}
	
	fullPrompt := fmt.Sprintf("%s\n\nContent to analyze:\n%s", prompt, content)
	
	reqBody := localRequest{
		Model:       l.config.Model,
		MaxTokens:   l.config.MaxTokens,
		Temperature: l.config.Temperature,
		Stream:      false,
		Messages: []message{
			{
				Role:    "user",
				Content: fullPrompt,
			},
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("Failed to prepare request: %v", err), "local")
	}
	
	endpoint := fmt.Sprintf("%s/v1/chat/completions", l.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("Failed to create HTTP request: %v", err), "local")
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := l.client.Do(req)
	if err != nil {
		// Check for specific error types
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return nil, WrapTimeoutError(err, "local")
			}
			// For local provider, connection refused likely means service is down
			if strings.Contains(err.Error(), "connection refused") {
				return nil, NewLLMError(ErrorTypeService, fmt.Sprintf("Local LLM service not available at %s", l.config.BaseURL), "local")
			}
		}
		return nil, WrapNetworkError(err, "local")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, WrapNetworkError(fmt.Errorf("failed to read response body: %w", err), "local")
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, ParseHTTPError(resp.StatusCode, body, "local")
	}
	
	var localResp localResponse
	if err := json.Unmarshal(body, &localResp); err != nil {
		return nil, WrapResponseError(fmt.Errorf("failed to parse response JSON: %w", err), "local")
	}
	
	if len(localResp.Choices) == 0 {
		return nil, NewLLMError(ErrorTypeResponse, "No choices in local LLM response", "local")
	}
	
	if localResp.Choices[0].Message.Content == "" {
		return nil, NewLLMError(ErrorTypeResponse, "Empty message content in local LLM response", "local")
	}
	
	return &Response{
		Content:    localResp.Choices[0].Message.Content,
		TokensUsed: localResp.Usage.TotalTokens,
		Model:      localResp.Model,
		Metadata: map[string]any{
			"prompt_tokens":     localResp.Usage.PromptTokens,
			"completion_tokens": localResp.Usage.CompletionTokens,
		},
	}, nil
}