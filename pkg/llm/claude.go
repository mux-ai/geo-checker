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

type ClaudeProvider struct {
	config *ProviderConfig
	client *http.Client
}

type claudeRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Messages    []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Model string `json:"model"`
}

func NewClaudeProvider(config *ProviderConfig) (*ClaudeProvider, error) {
	if config == nil {
		return nil, NewLLMError(ErrorTypeRequest, "Provider configuration is required", "claude")
	}
	
	if config.APIKey == "" {
		return nil, NewLLMError(ErrorTypeAuth, "Claude API key is required (set CLAUDE_API_KEY environment variable)", "claude")
	}
	
	// Validate API key format (should start with 'sk-ant-')
	if !strings.HasPrefix(config.APIKey, "sk-ant-") {
		return nil, NewLLMError(ErrorTypeAuth, "Invalid Claude API key format (should start with 'sk-ant-')", "claude")
	}
	
	if config.Model == "" {
		config.Model = "claude-3-sonnet-20240229"
	}
	
	// Validate model name
	validModels := []string{
		"claude-3-sonnet-20240229",
		"claude-3-opus-20240229",
		"claude-3-haiku-20240307",
		"claude-3-5-sonnet-20241022",
	}
	isValidModel := false
	for _, validModel := range validModels {
		if config.Model == validModel {
			isValidModel = true
			break
		}
	}
	if !isValidModel {
		return nil, NewLLMError(ErrorTypeModel, fmt.Sprintf("Unsupported Claude model: %s", config.Model), "claude")
	}
	
	if config.MaxTokens == 0 {
		config.MaxTokens = 4000
	}
	
	// Validate token limits
	if config.MaxTokens < 1 || config.MaxTokens > 8192 {
		return nil, NewLLMError(ErrorTypeRequest, "MaxTokens must be between 1 and 8192 for Claude", "claude")
	}
	
	// Validate temperature
	if config.Temperature < 0 || config.Temperature > 1 {
		return nil, NewLLMError(ErrorTypeRequest, "Temperature must be between 0 and 1", "claude")
	}
	
	return &ClaudeProvider{
		config: config,
		client: &http.Client{Timeout: 60 * time.Second},
	}, nil
}

func (c *ClaudeProvider) Name() string {
	return "claude"
}

func (c *ClaudeProvider) Analyze(ctx context.Context, content string, prompt string) (*Response, error) {
	// Validate inputs
	if strings.TrimSpace(content) == "" {
		return nil, NewLLMError(ErrorTypeRequest, "Content cannot be empty - webpage scraping may have failed or returned no extractable content", "claude")
	}
	if strings.TrimSpace(prompt) == "" {
		return nil, NewLLMError(ErrorTypeRequest, "Prompt cannot be empty", "claude")
	}
	
	fullPrompt := fmt.Sprintf("%s\n\nContent to analyze:\n%s", prompt, content)
	
	// Check content length
	if len(fullPrompt) > 200000 { // Claude's approximate context limit
		return nil, NewLLMError(ErrorTypeRequest, "Content too long for Claude model", "claude")
	}
	
	reqBody := claudeRequest{
		Model:       c.config.Model,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		Messages: []message{
			{
				Role:    "user",
				Content: fullPrompt,
			},
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("Failed to prepare request: %v", err), "claude")
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("Failed to create HTTP request: %v", err), "claude")
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	resp, err := c.client.Do(req)
	if err != nil {
		// Check for specific error types
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return nil, WrapTimeoutError(err, "claude")
			}
		}
		return nil, WrapNetworkError(err, "claude")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, WrapNetworkError(fmt.Errorf("failed to read response body: %w", err), "claude")
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, ParseHTTPError(resp.StatusCode, body, "claude")
	}
	
	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, WrapResponseError(fmt.Errorf("failed to parse response JSON: %w", err), "claude")
	}
	
	if len(claudeResp.Content) == 0 {
		return nil, NewLLMError(ErrorTypeResponse, "No content in Claude response", "claude")
	}
	
	if claudeResp.Content[0].Text == "" {
		return nil, NewLLMError(ErrorTypeResponse, "Empty text content in Claude response", "claude")
	}
	
	return &Response{
		Content:    claudeResp.Content[0].Text,
		TokensUsed: claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		Model:      claudeResp.Model,
		Metadata: map[string]any{
			"input_tokens":  claudeResp.Usage.InputTokens,
			"output_tokens": claudeResp.Usage.OutputTokens,
		},
	}, nil
}