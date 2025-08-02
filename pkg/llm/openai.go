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

type OpenAIProvider struct {
	config *ProviderConfig
	client *http.Client
}

type openAIRequest struct {
	Model       string        `json:"model"`
	Messages    []message     `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type openAIResponse struct {
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

func NewOpenAIProvider(config *ProviderConfig) (*OpenAIProvider, error) {
	if config == nil {
		return nil, NewLLMError(ErrorTypeRequest, "Provider configuration is required", "openai")
	}
	
	if config.APIKey == "" {
		return nil, NewLLMError(ErrorTypeAuth, "OpenAI API key is required (set OPENAI_API_KEY environment variable)", "openai")
	}
	
	// Validate API key format (should start with 'sk-')
	if !strings.HasPrefix(config.APIKey, "sk-") {
		return nil, NewLLMError(ErrorTypeAuth, "Invalid OpenAI API key format (should start with 'sk-')", "openai")
	}
	
	if config.Model == "" {
		config.Model = "gpt-4"
	}
	
	// Validate model name
	validModels := []string{
		"gpt-4", "gpt-4-turbo", "gpt-4-turbo-preview",
		"gpt-3.5-turbo", "gpt-3.5-turbo-16k",
		"gpt-4o", "gpt-4o-mini",
	}
	isValidModel := false
	for _, validModel := range validModels {
		if config.Model == validModel {
			isValidModel = true
			break
		}
	}
	if !isValidModel {
		return nil, NewLLMError(ErrorTypeModel, fmt.Sprintf("Unsupported OpenAI model: %s", config.Model), "openai")
	}
	
	if config.MaxTokens == 0 {
		config.MaxTokens = 4000
	}
	
	// Validate token limits based on model
	maxAllowed := 4096
	if strings.Contains(config.Model, "16k") {
		maxAllowed = 16384
	} else if strings.Contains(config.Model, "turbo") || strings.Contains(config.Model, "4o") {
		maxAllowed = 8192
	}
	
	if config.MaxTokens < 1 || config.MaxTokens > maxAllowed {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("MaxTokens must be between 1 and %d for model %s", maxAllowed, config.Model), "openai")
	}
	
	// Validate temperature
	if config.Temperature < 0 || config.Temperature > 2 {
		return nil, NewLLMError(ErrorTypeRequest, "Temperature must be between 0 and 2 for OpenAI", "openai")
	}
	
	return &OpenAIProvider{
		config: config,
		client: &http.Client{Timeout: 60 * time.Second},
	}, nil
}

func (o *OpenAIProvider) Name() string {
	return "openai"
}

func (o *OpenAIProvider) Analyze(ctx context.Context, content string, prompt string) (*Response, error) {
	// Validate inputs
	if strings.TrimSpace(content) == "" {
		return nil, NewLLMError(ErrorTypeRequest, "Content cannot be empty - webpage scraping may have failed or returned no extractable content", "openai")
	}
	if strings.TrimSpace(prompt) == "" {
		return nil, NewLLMError(ErrorTypeRequest, "Prompt cannot be empty", "openai")
	}
	
	fullPrompt := fmt.Sprintf("%s\n\nContent to analyze:\n%s", prompt, content)
	
	// Check content length (approximate token count)
	if len(fullPrompt) > 100000 { // Rough estimate for token limits
		return nil, NewLLMError(ErrorTypeRequest, "Content too long for OpenAI model", "openai")
	}
	
	reqBody := openAIRequest{
		Model:       o.config.Model,
		MaxTokens:   o.config.MaxTokens,
		Temperature: o.config.Temperature,
		Messages: []message{
			{
				Role:    "user",
				Content: fullPrompt,
			},
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("Failed to prepare request: %v", err), "openai")
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, NewLLMError(ErrorTypeRequest, fmt.Sprintf("Failed to create HTTP request: %v", err), "openai")
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.config.APIKey)
	
	resp, err := o.client.Do(req)
	if err != nil {
		// Check for specific error types
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return nil, WrapTimeoutError(err, "openai")
			}
		}
		return nil, WrapNetworkError(err, "openai")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, WrapNetworkError(fmt.Errorf("failed to read response body: %w", err), "openai")
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, ParseHTTPError(resp.StatusCode, body, "openai")
	}
	
	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, WrapResponseError(fmt.Errorf("failed to parse response JSON: %w", err), "openai")
	}
	
	if len(openAIResp.Choices) == 0 {
		return nil, NewLLMError(ErrorTypeResponse, "No choices in OpenAI response", "openai")
	}
	
	if openAIResp.Choices[0].Message.Content == "" {
		return nil, NewLLMError(ErrorTypeResponse, "Empty message content in OpenAI response", "openai")
	}
	
	return &Response{
		Content:    openAIResp.Choices[0].Message.Content,
		TokensUsed: openAIResp.Usage.TotalTokens,
		Model:      openAIResp.Model,
		Metadata: map[string]any{
			"prompt_tokens":     openAIResp.Usage.PromptTokens,
			"completion_tokens": openAIResp.Usage.CompletionTokens,
		},
	}, nil
}