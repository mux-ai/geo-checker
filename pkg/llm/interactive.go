package llm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ModelInfo contains information about available models
type ModelInfo struct {
	Name        string
	Provider    string
	Description string
	MaxTokens   int
	Recommended bool
}

// GetAvailableModels returns a list of available models for each provider
func GetAvailableModels() map[string][]ModelInfo {
	return map[string][]ModelInfo{
		"claude": {
			{
				Name:        "claude-3-5-sonnet-20241022",
				Provider:    "claude",
				Description: "Latest Claude 3.5 Sonnet - Best for complex analysis",
				MaxTokens:   8192,
				Recommended: true,
			},
			{
				Name:        "claude-3-sonnet-20240229",
				Provider:    "claude",
				Description: "Claude 3 Sonnet - Balanced performance and cost",
				MaxTokens:   4096,
				Recommended: false,
			},
			{
				Name:        "claude-3-opus-20240229",
				Provider:    "claude",
				Description: "Claude 3 Opus - Most capable, higher cost",
				MaxTokens:   4096,
				Recommended: false,
			},
			{
				Name:        "claude-3-haiku-20240307",
				Provider:    "claude",
				Description: "Claude 3 Haiku - Fastest, most economical",
				MaxTokens:   4096,
				Recommended: false,
			},
		},
		"openai": {
			{
				Name:        "gpt-4o",
				Provider:    "openai",
				Description: "GPT-4 Omni - Latest multimodal model, best overall",
				MaxTokens:   8192,
				Recommended: true,
			},
			{
				Name:        "gpt-4o-mini",
				Provider:    "openai",
				Description: "GPT-4 Omni Mini - Cost-effective, fast, excellent value",
				MaxTokens:   8192,
				Recommended: false,
			},
			{
				Name:        "gpt-4-turbo",
				Provider:    "openai",
				Description: "GPT-4 Turbo - Large context window, strong reasoning",
				MaxTokens:   8192,
				Recommended: false,
			},
			{
				Name:        "gpt-4",
				Provider:    "openai",
				Description: "GPT-4 - High quality reasoning, proven performance",
				MaxTokens:   4096,
				Recommended: false,
			},
			{
				Name:        "gpt-3.5-turbo",
				Provider:    "openai",
				Description: "GPT-3.5 Turbo - Fast, economical, good for simple tasks",
				MaxTokens:   4096,
				Recommended: false,
			},
		},
		"local": {
			{
				Name:        "llama2",
				Provider:    "local",
				Description: "Llama 2 - Open source model",
				MaxTokens:   4096,
				Recommended: true,
			},
			{
				Name:        "llama3",
				Provider:    "local",
				Description: "Llama 3 - Latest open source model",
				MaxTokens:   8192,
				Recommended: false,
			},
			{
				Name:        "codellama",
				Provider:    "local",
				Description: "Code Llama - Specialized for code analysis",
				MaxTokens:   4096,
				Recommended: false,
			},
			{
				Name:        "mistral",
				Provider:    "local",
				Description: "Mistral - Efficient open source model",
				MaxTokens:   4096,
				Recommended: false,
			},
		},
	}
}

// InteractiveModelSelection provides an interactive CLI for model selection
func InteractiveModelSelection(currentProvider string) (string, string, error) {
	models := GetAvailableModels()
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Select provider if not specified or confirm current
	var selectedProvider string
	if currentProvider == "" {
		fmt.Println("🤖 LLM Provider Selection")
		fmt.Println("========================")
		fmt.Println("1. claude   - Anthropic Claude models (requires CLAUDE_API_KEY)")
		fmt.Println("2. openai   - OpenAI GPT models (requires OPENAI_API_KEY)")
		fmt.Println("3. local    - Local LLM server (requires local server running)")
		fmt.Println()
		fmt.Print("Select provider (1-3): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", "", fmt.Errorf("failed to read input: %w", err)
		}

		choice := strings.TrimSpace(input)
		switch choice {
		case "1":
			selectedProvider = "claude"
		case "2":
			selectedProvider = "openai"
		case "3":
			selectedProvider = "local"
		default:
			return "", "", fmt.Errorf("invalid choice: %s", choice)
		}
	} else {
		selectedProvider = currentProvider
		fmt.Printf("Using provider: %s\n", selectedProvider)
	}

	// Step 2: Select model for the chosen provider
	providerModels, exists := models[selectedProvider]
	if !exists {
		return "", "", fmt.Errorf("no models available for provider: %s", selectedProvider)
	}

	fmt.Printf("\n📋 Available %s models:\n", strings.ToUpper(selectedProvider))
	fmt.Println(strings.Repeat("=", 50))

	for i, model := range providerModels {
		indicator := " "
		if model.Recommended {
			indicator = "⭐"
		}
		fmt.Printf("%s %d. %s\n", indicator, i+1, model.Name)
		fmt.Printf("    %s\n", model.Description)
		fmt.Printf("    Max tokens: %d\n", model.MaxTokens)
		fmt.Println()
	}

	fmt.Print("Select model (1-" + strconv.Itoa(len(providerModels)) + "): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("failed to read input: %w", err)
	}

	choiceNum, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choiceNum < 1 || choiceNum > len(providerModels) {
		return "", "", fmt.Errorf("invalid choice: %s", strings.TrimSpace(input))
	}

	selectedModel := providerModels[choiceNum-1]
	
	fmt.Printf("\n✅ Selected: %s (%s)\n", selectedModel.Name, selectedModel.Description)

	return selectedProvider, selectedModel.Name, nil
}

// ValidateModelForProvider checks if a model is valid for the given provider
func ValidateModelForProvider(provider, model string) error {
	models := GetAvailableModels()
	providerModels, exists := models[provider]
	if !exists {
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	for _, m := range providerModels {
		if m.Name == model {
			return nil
		}
	}

	// For local provider, be more permissive as users might have custom models
	if provider == "local" {
		return nil
	}

	return fmt.Errorf("unsupported model '%s' for provider '%s'", model, provider)
}

// GetRecommendedModel returns the recommended model for a provider
func GetRecommendedModel(provider string) string {
	models := GetAvailableModels()
	providerModels, exists := models[provider]
	if !exists {
		return ""
	}

	for _, model := range providerModels {
		if model.Recommended {
			return model.Name
		}
	}

	// Fallback to first model if no recommended model
	if len(providerModels) > 0 {
		return providerModels[0].Name
	}

	return ""
}

