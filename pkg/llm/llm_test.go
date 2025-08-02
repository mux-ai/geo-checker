package llm

import (
	"context"
	"testing"
)

func TestNewClaudeProvider_InvalidAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		config  *ProviderConfig
		wantErr string
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: "Provider configuration is required",
		},
		{
			name: "empty API key",
			config: &ProviderConfig{
				APIKey: "",
			},
			wantErr: "Claude API key is required",
		},
		{
			name: "invalid API key format",
			config: &ProviderConfig{
				APIKey: "invalid-key",
			},
			wantErr: "Invalid Claude API key format",
		},
		{
			name: "invalid model",
			config: &ProviderConfig{
				APIKey: "sk-ant-valid-format",
				Model:  "invalid-model",
			},
			wantErr: "Unsupported Claude model",
		},
		{
			name: "invalid max tokens",
			config: &ProviderConfig{
				APIKey:    "sk-ant-valid-format",
				Model:     "claude-3-sonnet-20240229",
				MaxTokens: 10000,
			},
			wantErr: "MaxTokens must be between 1 and 8192",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClaudeProvider(tt.config)
			if err == nil {
				t.Errorf("NewClaudeProvider() expected error, got nil")
				return
			}
			
			if llmErr, ok := err.(*LLMError); ok {
				if llmErr.Message != tt.wantErr && !contains(llmErr.Message, tt.wantErr) {
					t.Errorf("NewClaudeProvider() error = %v, wantErr %v", llmErr.Message, tt.wantErr)
				}
			} else {
				t.Errorf("NewClaudeProvider() expected LLMError, got %T", err)
			}
		})
	}
}

func TestNewOpenAIProvider_InvalidAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		config  *ProviderConfig
		wantErr string
	}{
		{
			name: "invalid API key format",
			config: &ProviderConfig{
				APIKey: "invalid-key",
			},
			wantErr: "Invalid OpenAI API key format",
		},
		{
			name: "invalid model",
			config: &ProviderConfig{
				APIKey: "sk-valid-format",
				Model:  "invalid-model",
			},
			wantErr: "Unsupported OpenAI model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOpenAIProvider(tt.config)
			if err == nil {
				t.Errorf("NewOpenAIProvider() expected error, got nil")
				return
			}
		})
	}
}

func TestAnalyze_InvalidInputs(t *testing.T) {
	config := &ProviderConfig{
		APIKey: "sk-ant-test-key",
		Model:  "claude-3-sonnet-20240229",
	}
	
	provider, err := NewClaudeProvider(config)
	if err != nil {
		t.Fatalf("NewClaudeProvider() failed: %v", err)
	}

	tests := []struct {
		name    string
		content string
		prompt  string
		wantErr string
	}{
		{
			name:    "empty content",
			content: "",
			prompt:  "test prompt",
			wantErr: "Content cannot be empty",
		},
		{
			name:    "empty prompt",
			content: "test content",
			prompt:  "",
			wantErr: "Prompt cannot be empty",
		},
		{
			name:    "content too long",
			content: string(make([]byte, 300000)), // Very long content
			prompt:  "test prompt",
			wantErr: "Content too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := provider.Analyze(ctx, tt.content, tt.prompt)
			if err == nil {
				t.Errorf("Analyze() expected error, got nil")
				return
			}
			
			if llmErr, ok := err.(*LLMError); ok {
				if !contains(llmErr.Message, tt.wantErr) {
					t.Errorf("Analyze() error = %v, wantErr %v", llmErr.Message, tt.wantErr)
				}
			}
		})
	}
}

func TestGetRecommendedModel(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		want     string
	}{
		{
			name:     "claude",
			provider: "claude",
			want:     "claude-3-5-sonnet-20241022",
		},
		{
			name:     "openai",
			provider: "openai",
			want:     "gpt-4o",
		},
		{
			name:     "local",
			provider: "local",
			want:     "llama2",
		},
		{
			name:     "unknown provider",
			provider: "unknown",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRecommendedModel(tt.provider)
			if got != tt.want {
				t.Errorf("GetRecommendedModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateModelForProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		model    string
		wantErr  bool
	}{
		{
			name:     "valid claude model",
			provider: "claude",
			model:    "claude-3-sonnet-20240229",
			wantErr:  false,
		},
		{
			name:     "invalid claude model",
			provider: "claude",
			model:    "invalid-model",
			wantErr:  true,
		},
		{
			name:     "valid openai model",
			provider: "openai",
			model:    "gpt-4",
			wantErr:  false,
		},
		{
			name:     "local provider with any model",
			provider: "local",
			model:    "custom-model",
			wantErr:  false, // Local provider is permissive
		},
		{
			name:     "unknown provider",
			provider: "unknown",
			model:    "any-model",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModelForProvider(tt.provider, tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModelForProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
			 s[len(s)-len(substr):] == substr || 
			 containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}